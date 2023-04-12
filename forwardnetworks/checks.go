package forwardnetworks

import (
	"encoding/json"
	"errors"
	"fmt"
	"bytes"
	"net/http"
	"net/url"
	"io/ioutil"
)

type Check struct {
	ID                    string `json:"id"`
	Definition            CheckDefinition `json:"definition"`
	Enabled               bool `json:"enabled"`
	Priority              string `json:"priority"`
	Name                  string `json:"name"`
	CreationDateMillis    int64 `json:"creationDateMillis"`
	CreatorId             string `json:"creatorId"`
	DefinitionDateMillis  int64 `json:"definitionDateMillis"`
	Description           string `json:"description"`
	Status                string `json:"status"`
	ExecutionDateMillis   int64 `json:"executionDateMillis"`
	ExecutionDurationMillis int64 `json:"executionDurationMillis"`
}

type CheckDefinition struct {
	PredefinedCheckType string `json:"predefinedCheckType"`
	CheckType           string `json:"checkType"`
	QueryID             string `json:"queryId,omitempty"`
}

// GetChecks retrieves the status of all checks.
func (c *ForwardNetworksClient) GetChecks(snapshotID string, checkType, priority, status string) ([]Check, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/snapshots/%s/checks", c.BaseURL, snapshotID))
	if err != nil {
		return nil, err
	}

	query := baseURL.Query()
	if checkType != "" {
		query.Set("type", checkType)
	}
	if priority != "" {
		query.Set("priority", priority)
	}
	if status != "" {
		query.Set("status", status)
	}
	baseURL.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusConflict {
		var apiError ApiError
		err = json.Unmarshal(bodyBytes, &apiError)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error getting checks")
	}

	var checks []Check
	err = json.Unmarshal(bodyBytes, &checks)
	if err != nil {
		return nil, err
	}

	return checks, nil
}

func (c *ForwardNetworksClient) GetCheck(snapshotID string, checkID string) (*Check, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/snapshots/%s/checks/%s", c.BaseURL, snapshotID, checkID))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusConflict {
		var apiError ApiError
		err = json.Unmarshal(bodyBytes, &apiError)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error getting check")
	}

	var check Check
	err = json.Unmarshal(bodyBytes, &check)
	if err != nil {
		return nil, err
	}

	return &check, nil
}

func (c *ForwardNetworksClient) ActivateCheck(snapshotID string, checkType, queryID string, enabled bool) (*Check, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/snapshots/%s/checks", c.BaseURL, snapshotID))
	if err != nil {
		return nil, err
	}

	newCheck := map[string]interface{}{
		"definition": map[string]string{
			"checkType": checkType,
			"queryId":   queryID,
		},
		"enabled": enabled,
	}

	newCheckJSON, err := json.Marshal(newCheck)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL.String(), bytes.NewBuffer(newCheckJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusConflict {
		var apiError ApiError
		err = json.Unmarshal(bodyBytes, &apiError)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New("error activating check")
	}

	var check Check
	err = json.Unmarshal(bodyBytes, &check)
	if err != nil {
		return nil, err
	}

	return &check, nil
}

func (c *ForwardNetworksClient) DeactivateCheck(snapshotID string, checkID string) error {
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/snapshots/%s/checks/%s", c.BaseURL, snapshotID, checkID))
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, baseURL.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var apiError ApiError
		err = json.Unmarshal(bodyBytes, &apiError)
		if err != nil {
			return err
		}
		return apiError
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New("error deactivating check")
	}

	return nil
}

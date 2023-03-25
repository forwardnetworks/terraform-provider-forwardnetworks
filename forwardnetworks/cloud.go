package forwardnetworks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CloudAccount struct {
	Name                  string                 `json:"name"`
	Type                  string                 `json:"type"`
	Collect               bool                   `json:"collect"`
	NumVirtualizedDevices int                    `json:"numVirtualizedDevices"`
	Regions               map[string]interface{} `json:"regions,omitempty"`
	AssumeRoleInfos       []interface{}          `json:"assumeRoleInfos,omitempty"`
	Subscriptions         []interface{}          `json:"subscriptions,omitempty"`
}

func (c *ForwardNetworksClient) ListCloudAccounts(networkID string) ([]*CloudAccount, error) {
	url := fmt.Sprintf("%s/api/networks/%s/cloudAccounts", c.BaseURL, networkID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cloudAccounts []*CloudAccount
	err = json.Unmarshal(bodyBytes, &cloudAccounts)
	if err != nil {
		return nil, err
	}

	return cloudAccounts, nil
}

func (c *ForwardNetworksClient) GetCloudAccount(networkID, setupID string) (*CloudAccount, error) {
	url := fmt.Sprintf("%s/api/networks/%s/cloudAccounts/%s", c.BaseURL, networkID, setupID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cloudAccount CloudAccount
	err = json.Unmarshal(bodyBytes, &cloudAccount)
	if err != nil {
		return nil, err
	}

	return &cloudAccount, nil
}




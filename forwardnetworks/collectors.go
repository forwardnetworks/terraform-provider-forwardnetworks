package forwardnetworks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
)

type Collector struct {
	OrgID             string `json:"orgId"`
	UserID            string `json:"userId"`
	Username          string `json:"username"`
	CollectorName     string `json:"collectorName"`
	ID                string `json:"id"`
	FWCollectorPoolAccount bool `json:"fwCollectorPoolAccount"`
	Status            struct {
		BusyStatus           string `json:"busyStatus"`
		Outdated             bool   `json:"outdated"`
		SupportsRemoteUpgrade bool   `json:"supportsRemoteUpgrade"`
	} `json:"status"`
}

type CollectorUser struct {
	Username string `json:"username"`
}

func (c *ForwardNetworksClient) GetCollectors() ([]Collector, error) {
	url := fmt.Sprintf("%s/api/orgs/current/collectors", c.BaseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d or URL: %s", resp.StatusCode, url)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var collectors []Collector
	err = json.Unmarshal(bodyBytes, &collectors)
	if err != nil {
		return nil, err
	}

	return collectors, nil
}

func (c *ForwardNetworksClient) GetCollector(collectorID string) (*Collector, error) {
	collectors, err := c.GetCollectors()
	if err != nil {
		return nil, err
	}

	for _, collector := range collectors {
		if collector.ID == collectorID {
			return &collector, nil
		}
	}

	return nil, fmt.Errorf("collector with ID %s not found", collectorID)
}

func (c *ForwardNetworksClient) UpdateCollector(networkID, collectorName, collectorUsername string) error {
    url := fmt.Sprintf("%s/api/networks/%s/collector/user", c.BaseURL, networkID)

    collectorUser := &CollectorUser{
        Username: collectorUsername, // Update this line
    }

    collectorUserData, err := json.Marshal(collectorUser)
    if err != nil {
        return err
    }

    req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(collectorUserData))
    if err != nil {
        return err
    }

    req.SetBasicAuth(c.Username, c.Password)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("API request failed with status code %d or URL: %s", resp.StatusCode, url)
    }

    return nil
}
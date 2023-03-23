package forwardnetworks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ForwardNetworksClient struct {
	Username string
	Password string
	BaseURL  string
}

func NewForwardNetworksClient(username, password, baseURL string) *ForwardNetworksClient {
	return &ForwardNetworksClient{
		Username: username,
		Password: password,
		BaseURL:  baseURL,
	}
}

func (c *ForwardNetworksClient) GetVersion() (string, error) {
	url := fmt.Sprintf("%s/api/version", c.BaseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(c.Username, c.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var versionData struct {
		Version string `json:"version"`
	}
	err = json.Unmarshal(bodyBytes, &versionData)
	if err != nil {
		return "", err
	}

	return versionData.Version, nil
}


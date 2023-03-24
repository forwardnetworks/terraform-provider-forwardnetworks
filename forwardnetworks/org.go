package forwardnetworks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *ForwardNetworksClient) GetOrgId() (string, string, error) {
	url := fmt.Sprintf("%s/api/orgs/current", c.BaseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}

	req.SetBasicAuth(c.Username, c.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var orgData struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	err = json.Unmarshal(bodyBytes, &orgData)
	if err != nil {
		return "", "", err
	}

	return orgData.Id, orgData.Name, nil
}
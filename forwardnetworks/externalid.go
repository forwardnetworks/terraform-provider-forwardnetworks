package forwardnetworks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *ForwardNetworksClient) GetExternalId(networkId string) (string, error) {
	url := fmt.Sprintf("%s/api/networks/%s/cloudAccounts/aws/assumeRole/externalId", c.BaseURL, networkId)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HttpClient.Do(req)
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

	var externalIdData struct {
		ExternalId string `json:"externalId"`
	}
	err = json.Unmarshal(bodyBytes, &externalIdData)
	if err != nil {
		return "", err
	}

	return externalIdData.ExternalId, nil
}

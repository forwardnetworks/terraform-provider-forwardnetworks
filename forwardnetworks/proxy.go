package forwardnetworks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Proxy struct {
	Id                  string `json:"id"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	Username            string `json:"username"`
	Protocol            string `json:"protocol"`
	DisableCertChecking bool   `json:"disableCertChecking"`
}

func (c *ForwardNetworksClient) GetProxy(networkId string) (*Proxy, error) {
	url := fmt.Sprintf("%s/api/networks/%s/proxy", c.BaseURL, networkId)

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

var proxy Proxy
err = json.Unmarshal(bodyBytes, &proxy)
if err != nil {
	return nil, err
}

return &proxy, nil
}
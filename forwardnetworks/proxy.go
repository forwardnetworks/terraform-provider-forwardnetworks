package forwardnetworks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
)

type Proxy struct {
	Id                  string `json:"id,omitempty"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	Username            string `json:"username"`
	Password            string `json:"password"`
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

var proxy Proxy
err = json.Unmarshal(bodyBytes, &proxy)
if err != nil {
	return nil, err
}

return &proxy, nil
}

func (c *ForwardNetworksClient) CreateOrUpdateProxy(networkId string, proxy *Proxy) error {
	url := fmt.Sprintf("%s/api/networks/%s/proxy", c.BaseURL, networkId)

	proxyData, err := json.Marshal(proxy)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(proxyData))
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
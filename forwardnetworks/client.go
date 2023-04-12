package forwardnetworks

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ForwardNetworksClient struct {
	Username string
	Password string
	BaseURL  string
	Insecure bool
	HttpClient *http.Client
}

func NewForwardNetworksClient(baseURL, username, password string, insecure bool) *ForwardNetworksClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}

	return &ForwardNetworksClient{
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		Insecure:   insecure,
		HttpClient: httpClient,
	}
}


func (c *ForwardNetworksClient) GetVersion() (string, error) {
	url := fmt.Sprintf("%s/api/version", c.BaseURL)

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

	var versionData struct {
		Version string `json:"version"`
	}
	err = json.Unmarshal(bodyBytes, &versionData)
	if err != nil {
		return "", err
	}

	return versionData.Version, nil
}


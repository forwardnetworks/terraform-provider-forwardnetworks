package forwardnetworks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AwsPolicy struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

type Statement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource string   `json:"Resource"`
}

func (c *ForwardNetworksClient) GetAwsPolicy() (*AwsPolicy, error) {
	url := fmt.Sprintf("%s/api/cloud/aws-policy", c.BaseURL)

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
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var awsPolicy AwsPolicy
	err = json.Unmarshal(bodyBytes, &awsPolicy)
	if err != nil {
		return nil, err
	}

	return &awsPolicy, nil
}


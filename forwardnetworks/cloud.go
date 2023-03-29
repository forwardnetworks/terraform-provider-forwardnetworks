package forwardnetworks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CloudAccount struct {
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`
	Collect         bool                   `json:"collect"`
	Regions         map[string]interface{}	`json:"regions,omitempty"`
	AssumeRoleInfos []AssumeRoleInfo       `json:"assumeRoleInfos,omitempty"`
	Subscriptions   []Subscription         `json:"subscriptions,omitempty"`
	AWSUsername     string                 `json:"username,omitempty"`
    ProxyServerId     string                `json:"proxyServerId,omitempty"`
}

type AssumeRoleInfo struct {
	AccountId   string `json:"accountId"`
	AccountName string `json:"accountName"`
	RoleArn     string `json:"rolearn"`
	Enabled     bool   `json:"enabled"`
}

type Subscription struct {
	SubscriptionId string `json:"subscriptionId"`
	ClientId       string `json:"clientId"`
	Tenant         string `json:"tenant"`
	Enabled        bool   `json:"enabled"`
}

func (c *ForwardNetworksClient) GetCloudAccounts(networkId string, accountName string) (map[string]CloudAccount, error) {
    url := fmt.Sprintf("%s/api/networks/%s/cloudAccounts", c.BaseURL, networkId)
    
    cloudAccounts := make(map[string]CloudAccount)

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

    var cloudAccountsSlice []CloudAccount
    err = json.Unmarshal(bodyBytes, &cloudAccountsSlice)
    if err != nil {
        return nil, err
    }

    for _, cloudAccount := range cloudAccountsSlice {
        accountNames := []string{}
        if cloudAccount.Type == "AWS" {
            for _, assumeRoleInfo := range cloudAccount.AssumeRoleInfos {
                accountNames = append(accountNames, assumeRoleInfo.AccountName)
            }
        } else if cloudAccount.Type == "AZURE" {
            for _, subscription := range cloudAccount.Subscriptions {
                accountNames = append(accountNames, subscription.SubscriptionId)
            }

            for _, accountName := range accountNames {
                credentialUrl := fmt.Sprintf("%s/api/networks/%s/cloudAccounts/%s/credential", c.BaseURL, networkId, cloudAccount.Name)

                credentialReq, err := http.NewRequest(http.MethodGet, credentialUrl, nil)
                if err != nil {
                    return nil, err
                }

                credentialReq.SetBasicAuth(c.Username, c.Password)

                credentialResp, err := client.Do(credentialReq)
                if err != nil {
                    return nil, err
                }
                defer credentialResp.Body.Close()

                if credentialResp.StatusCode != http.StatusOK {
                    return nil, fmt.Errorf("API request failed with status code %d", credentialResp.StatusCode)
                }

                var credential struct {
                    Subscriptions []Subscription `json:"subscriptions"`
                }
                err = json.NewDecoder(credentialResp.Body).Decode(&credential)
                if err != nil {
                    return nil, err
                }

                for j, subscription := range credential.Subscriptions {
                    if subscription.SubscriptionId == accountName {
                        cloudAccount.Subscriptions[j].ClientId = subscription.ClientId
                        cloudAccount.Subscriptions[j].Tenant = subscription.Tenant
                        break
                    }
                }
            }
        } // Add other cloud providers here

        // extract desired values from regions field
        regions := make(map[string]interface{})
        for k, v := range cloudAccount.Regions {
            if r, ok := v.(map[string]interface{}); ok {
                if t, ok := r["testInstant"].(float64); ok {
                    regions[k] = int(t)
                }
            }
        }
        cloudAccount.Regions = regions

        // Add the cloudAccount to the map using the account name as the key
        cloudAccounts[cloudAccount.Name] = cloudAccount
    }
    	if accountName != "" {
			if account, ok := cloudAccounts[accountName]; ok {
				cloudAccounts = map[string]CloudAccount{
					accountName: account,
				}
		} else {
			return nil, fmt.Errorf("account with name %s not found", accountName)
		}
	}
    return cloudAccounts, nil
}

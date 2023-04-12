package forwardnetworks

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "bytes"
)

type Network struct {
    ID        string `json:"id"`
    ParentID  string `json:"parentId,omitempty"`
    Name      string `json:"name"`
    OrgID     string `json:"orgId"`
    Creator   string `json:"creator,omitempty"`
    CreatorID string `json:"creatorId,omitempty"`
    CreatedAt int64  `json:"createdAt"`
    Note      string `json:"note,omitempty"`
}

type NetworkUpdate struct {
    Name string `json:"name,omitempty"`
    Note string `json:"note,omitempty"`
}

type CreateWorkspaceNetworkRequest struct {
    Name         string   `json:"name"`
    Note         string   `json:"note,omitempty"`
    Devices      []string `json:"devices,omitempty"`
    CloudAccounts []string `json:"cloudAccounts,omitempty"`
    VCenters      []string `json:"vcenters,omitempty"`
}

func (c *ForwardNetworksClient) GetNetworks() ([]Network, error) {
    url := fmt.Sprintf("%s/api/networks", c.BaseURL)

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

    var networks []Network
    err = json.Unmarshal(bodyBytes, &networks)
    if err != nil {
        return nil, err
    }

    return networks, nil
}

func (c *ForwardNetworksClient) GetNetwork(networkId string) (*Network, error) {
    networks, err := c.GetNetworks()
    if err != nil {
        return nil, err
    }

    for _, network := range networks {
        if network.ID == networkId {
            return &network, nil
        }
    }

    return nil, fmt.Errorf("Network not found")
}

func (c *ForwardNetworksClient) CreateNetwork(newNetwork *Network) (*Network, error) {
    url := fmt.Sprintf("%s/api/networks?name=%s", c.BaseURL, newNetwork.Name)

    req, err := http.NewRequest(http.MethodPost, url, nil)
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

    var network Network
    err = json.NewDecoder(resp.Body).Decode(&network)
    if err != nil {
        return nil, err
    }

    return &network, nil
}

func (c *ForwardNetworksClient) UpdateNetwork(networkId string, update NetworkUpdate) (*Network, error) {
    url := fmt.Sprintf("%s/api/networks/%s", c.BaseURL, networkId)

    updateBytes, err := json.Marshal(update)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(updateBytes))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.SetBasicAuth(c.Username, c.Password)

    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code %d or URL: %s", resp.StatusCode, url)
    }

    var updatedNetwork Network
    err = json.NewDecoder(resp.Body).Decode(&updatedNetwork)
    if err != nil {
        return nil, err
    }

    return &updatedNetwork, nil
}

func (c *ForwardNetworksClient) DeleteNetwork(networkId string) error {
    url := fmt.Sprintf("%s/api/networks/%s", c.BaseURL, networkId)

    req, err := http.NewRequest(http.MethodDelete, url, nil)
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.SetBasicAuth(c.Username, c.Password)

    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
        return fmt.Errorf("API request failed with status code %d or URL: %s", resp.StatusCode, url)
    }

    return nil
}
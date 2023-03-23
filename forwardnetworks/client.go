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

func (c *ForwardNetworksClient) GetVersion() (string,


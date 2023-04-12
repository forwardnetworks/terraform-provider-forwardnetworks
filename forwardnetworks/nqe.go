package forwardnetworks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
)

type NQEQueryRead struct {
	QueryID       	string `json:"queryId"`
	Path 			string `json:"path"`
    Intent     		string `json:"intent"`
    Repository  	string `json:"repository"`
}

type NQEQueryBody struct {
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Query       string                 `json:"query"`
	QueryID       string                 `json:"queryId"`
	QueryOptions NQEQueryOptions        `json:"queryOptions,omitempty"`
}

type NQEQueryOptions struct {
	Offset        int               `json:"offset,omitempty"`
	Limit         int               `json:"limit,omitempty"`
	SortBy        NQEQuerySortBy    `json:"sortBy,omitempty"`
	ColumnFilters []NQEColumnFilter `json:"columnFilters,omitempty"`
}

type NQEQuerySortBy struct {
	ColumnName string `json:"columnName"`
	Order      string `json:"order"`
}

type NQEColumnFilter struct {
	ColumnName string      `json:"columnName"`
	Value      interface{} `json:"value"`
}

type NQEQueryResult struct {
	APIUrl         string               `json:"apiUrl"`
	HTTPMethod     string               `json:"httpMethod"`
	Message        string               `json:"message"`
	Reason         string               `json:"reason"`
	Errors         []NQEQueryResultError `json:"errors"`
	SnapshotID     string               `json:"snapshotId"`
	CompletionType string               `json:"completionType"`
}

type NQEQueryResultError struct {
	Message string `json:"message"`
}

func (c *ForwardNetworksClient) GetNqeQueries(path string) ([]NQEQueryRead, error) {
	url := c.BaseURL + "/api/nqe/queries"
	if path != "" {
		url += "?dir=" + path
	}

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

	var nqeQueries []NQEQueryRead
	err = json.Unmarshal(bodyBytes, &nqeQueries)
	if err != nil {
		return nil, err
	}

	return nqeQueries, nil
}

func (c *ForwardNetworksClient) ExecuteNQEQuery(networkID, snapshotID, query string, queryOptions map[string]interface{}, queryID ...string) ([]byte, error) {
    url := fmt.Sprintf("%s/api/nqe", c.BaseURL)

    if networkID != "" && snapshotID != "" {
        return nil, fmt.Errorf("only one of networkID and snapshotID should be supplied")
    }

    if networkID != "" {
        url += fmt.Sprintf("?networkId=%s", networkID)
    } else if snapshotID != "" {
        url += fmt.Sprintf("?snapshotId=%s", snapshotID)
    } else {
        return nil, fmt.Errorf("either networkID or snapshotID must be supplied")
    }

    requestBody := make(map[string]interface{})
    requestBody["query"] = query

    if len(queryID) > 0 {
        requestBody["queryId"] = queryID[0]
    }

    if len(queryOptions) > 0 {
        requestBody["queryOptions"] = queryOptions
    }

    requestBodyBytes, err := json.Marshal(requestBody)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyBytes))
    if err != nil {
        return nil, err
    }

    req.SetBasicAuth(c.Username, c.Password)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        errorResp := struct {
            APIUrl         string `json:"apiUrl"`
            HttpMethod     string `json:"httpMethod"`
            Message        string `json:"message"`
            Reason         string `json:"reason"`
            Errors         []struct {
                Message string `json:"message"`
            } `json:"errors"`
            SnapshotId     string `json:"snapshotId,omitempty"`
            CompletionType string `json:"completionType,omitempty"`
        }{}

        if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
            return nil, fmt.Errorf("API request failed with status code %d and invalid JSON response", resp.StatusCode)
        }

        if len(errorResp.Errors) > 0 {
            return nil, fmt.Errorf("query execution failed with errors: %v", errorResp.Errors)
        }

        return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, errorResp.Message)
    }

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return bodyBytes, nil
}

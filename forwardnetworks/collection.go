package forwardnetworks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CollectionSchedules struct {
	Schedules []CollectionSchedule `json:"schedules"`
}

type CollectionSchedule struct {
	Id             string   `json:"id"`
	Enabled        bool     `json:"enabled"`
	TimeZone       string   `json:"timeZone"`
	DaysOfTheWeek  []int    `json:"daysOfTheWeek"`
	Times          []string `json:"times,omitempty"`
	PeriodInSeconds int     `json:"periodInSeconds,omitempty"`
	StartAt        string   `json:"startAt,omitempty"`
	EndAt          string   `json:"endAt,omitempty"`
}

func (c *ForwardNetworksClient) GetCollectionSchedules(networkId string) (*CollectionSchedules, error) {
	url := fmt.Sprintf("%s/api/networks/%s/collection-schedules", c.BaseURL, networkId)

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

	var schedules CollectionSchedules
	err = json.Unmarshal(bodyBytes, &schedules)
	if err != nil {
		return nil, err
	}

	return &schedules, nil
}

func (c *ForwardNetworksClient) CreateCollectionSchedule(networkId string, schedule *CollectionSchedule) (*CollectionSchedule, error) {
	url := fmt.Sprintf("%s/api/networks/%s/collection-schedules", c.BaseURL, networkId)

	scheduleData, err := json.Marshal(schedule)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(scheduleData))
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
		return nil, fmt.Errorf("API request failed with status code %d or URL: %s", resp.StatusCode, url)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var createdSchedule CollectionSchedule
	err = json.Unmarshal(bodyBytes, &createdSchedule)
	if err != nil {
		return nil, err
	}

	return &createdSchedule, nil
}

func (c *ForwardNetworksClient) GetCollectionSchedule(networkID string, scheduleID string) (*CollectionSchedule, error) {
    url := fmt.Sprintf("%s/api/networks/%s/collection-schedules/%s", c.BaseURL, networkID, scheduleID)

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

    var collectionSchedule CollectionSchedule
    err = json.Unmarshal(bodyBytes, &collectionSchedule)
    if err != nil {
        return nil, err
    }

    return &collectionSchedule, nil
}

func (c *ForwardNetworksClient) UpdateCollectionSchedule(networkId string, scheduleId string, schedule *CollectionSchedule) (*CollectionSchedule, error) {
    url := fmt.Sprintf("%s/api/networks/%s/collections/schedules/%s", c.BaseURL, networkId, scheduleId)

    scheduleData, err := json.Marshal(schedule)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(scheduleData))
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
        return nil, fmt.Errorf("API request failed with status code %d or URL: %s", resp.StatusCode, url)
    }

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var updatedSchedule CollectionSchedule
    err = json.Unmarshal(bodyBytes, &updatedSchedule)
    if err != nil {
        return nil, err
    }

    return &updatedSchedule, nil
}

func (c *ForwardNetworksClient) DeleteCollectionSchedule(networkId string, scheduleId string) error {
	url := fmt.Sprintf("%s/api/networks/%s/collection-schedules/%s", c.BaseURL, networkId, scheduleId)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)

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

func (c *ForwardNetworksClient) StartCollection(networkId string) error {
	url := fmt.Sprintf("%s/api/networks/%s/startcollection", c.BaseURL, networkId)

	req, err := http.NewRequest(http.MethodPost, url, nil)
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

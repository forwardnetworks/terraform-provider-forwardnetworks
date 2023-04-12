package forwardnetworks

import (
   "context"
   "encoding/json"
   "fmt"
   "io/ioutil"
   "net/http"
   "bytes"
)

type Snapshot struct {
   ID                    string `json:"id"`
   ProcessingTrigger     string `json:"processingTrigger"`
   TotalDevices          int    `json:"totalDevices"`
   CreationDateMillis    int64  `json:"creationDateMillis"`
   ProcessedAtMillis     int64  `json:"processedAtMillis"`
   OldestCollectionMillis int64  `json:"oldestCollectionMillis"`
   LatestCollectionMillis int64  `json:"latestCollectionMillis"`
   IsDraft               bool   `json:"isDraft"`
   State                 string `json:"state"`
}

type Metrics struct {
   CollectionConcurrency          int                    `json:"collectionConcurrency"`
   CollectionDuration             int                    `json:"collectionDuration"`
   CollectionFailures             map[string]int         `json:"collectionFailures"`
   CreationDateMillis             int64                  `json:"creationDateMillis"`
   HostComputationStatus          string                 `json:"hostComputationStatus"`
   IpLocationIndexingStatus       string                 `json:"ipLocationIndexingStatus"`
   JumpServerCollectionConcurrency int                   `json:"jumpServerCollectionConcurrency"`
   L2IndexingStatus               string                 `json:"l2IndexingStatus"`
   NeedsReprocessing              bool                   `json:"needsReprocessing"`
   NumCollectionFailureDevices    int                    `json:"numCollectionFailureDevices"`
   NumParsingFailureDevices       int                    `json:"numParsingFailureDevices"`
   NumSuccessfulDevices           int                    `json:"numSuccessfulDevices"`
   ParsingFailures                map[string]int         `json:"parsingFailures"`
   PathSearchIndexingStatus       string                 `json:"pathSearchIndexingStatus"`
   ProcessingDuration             int                    `json:"processingDuration"`
   SearchIndexingStatus           string                 `json:"searchIndexingStatus"`
   SnapshotID                     string                 `json:"snapshotId"`
}

type ExportParams struct {
   IncludeDevices  []string `json:"includeDevices,omitempty"`
   ExcludeDevices  []string `json:"excludeDevices,omitempty"`
   ObfuscationKey  string   `json:"obfuscationKey,omitempty"`
   ObfuscateNames  bool     `json:"obfuscateNames,omitempty"`
}

type SnapshotWithMetrics struct {
   Snapshot
   Metrics *Metrics `json:"metrics,omitempty"`
}

type ApiError struct {
   ApiUrl     string `json:"apiUrl"`
   HttpMethod string `json:"httpMethod"`
   Message    string `json:"message"`
   Reason     string `json:"reason"`
}

func (e ApiError) Error() string {
    return fmt.Sprintf("API error: %s (%s %s) - reason: %s", e.Message, e.HttpMethod, e.ApiUrl, e.Reason)
}

func (c *ForwardNetworksClient) ListSnapshots(networkID string, latestProcessed, metrics bool) ([]SnapshotWithMetrics, error) {
    var url string

    if latestProcessed {
        url = fmt.Sprintf("%s/api/networks/%s/snapshots/latestProcessed", c.BaseURL, networkID)
    } else {
        url = fmt.Sprintf("%s/api/networks/%s/snapshots", c.BaseURL, networkID)
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

    if resp.StatusCode == http.StatusNotFound {
        return nil, ApiError{
            ApiUrl:     url,
            HttpMethod: http.MethodGet,
            Message:    "The network has no Snapshots.",
            Reason:     fmt.Sprintf("Status code: %d", resp.StatusCode),
        }
    } else if resp.StatusCode == http.StatusConflict {
        return nil, ApiError{
            ApiUrl:     url,
            HttpMethod: http.MethodGet,
            Message:    "None of the Snapshots in the network are processed. Processing of the latest Snapshot has begun.",
            Reason:     fmt.Sprintf("Status code: %d", resp.StatusCode),
        }
    } else if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code %d or URL: %s", resp.StatusCode, url)
    }

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if latestProcessed {
        var latestSnapshot Snapshot
        err = json.Unmarshal(bodyBytes, &latestSnapshot)
        if err != nil {
            return nil, err
        }
        snapshots := []Snapshot{latestSnapshot}
        if metrics {
            snapshotsWithMetrics, err := c.addMetricsToSnapshots(snapshots)
            if err != nil {
                return nil, err
            }
            return snapshotsWithMetrics, nil
        }
        return []SnapshotWithMetrics{{Snapshot: latestSnapshot}}, nil
    } else {
        var snapshotsResponse struct {
            Snapshots []Snapshot `json:"snapshots"`
        }
        err = json.Unmarshal(bodyBytes, &snapshotsResponse)
        if err != nil {
            return nil, err
        }

        if metrics {
            snapshotsWithMetrics, err := c.addMetricsToSnapshots(snapshotsResponse.Snapshots)
            if err != nil {
                return nil, err
            }
            return snapshotsWithMetrics, nil
        }
        snapshots := make([]SnapshotWithMetrics, len(snapshotsResponse.Snapshots))
        for i, snapshot := range snapshotsResponse.Snapshots {
            snapshots[i] = SnapshotWithMetrics{Snapshot: snapshot}
        }
        return snapshots, nil
    }
}


func (c *ForwardNetworksClient) addMetricsToSnapshots(snapshots []Snapshot) ([]SnapshotWithMetrics, error) {
   snapshotsWithMetrics := make([]SnapshotWithMetrics, len(snapshots))
      for i, snapshot := range snapshots {
      metrics, err := c.getSnapshotMetrics(snapshot.ID)
      if err != nil {
         return nil, err
      }
      snapshotsWithMetrics[i] = SnapshotWithMetrics{Snapshot: snapshot, Metrics: metrics}
   }
   return snapshotsWithMetrics, nil
}

func (c *ForwardNetworksClient) getSnapshotMetrics(snapshotID string) (*Metrics, error) {
   url := fmt.Sprintf("%s/api/snapshots/%s/metrics", c.BaseURL, snapshotID)

   req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
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
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var metrics Metrics
   err = json.NewDecoder(resp.Body).Decode(&metrics)
   if err != nil {
      return nil, err
   }

   return &metrics, nil
}

func (c *ForwardNetworksClient) ExportSnapshot(snapshotID string) ([]byte, error) {
   url := fmt.Sprintf("%s/api/snapshots/%s", c.BaseURL, snapshotID)

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

   data, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return data, nil
}

func (c *ForwardNetworksClient) ExportSnapshotWithParams(snapshotID string, params *ExportParams) ([]byte, error) {
   url := fmt.Sprintf("%s/api/snapshots/%s", c.BaseURL, snapshotID)

   paramsData, err := json.Marshal(params)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(paramsData))
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

   data, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return data, nil
}
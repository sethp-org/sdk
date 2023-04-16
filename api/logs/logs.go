package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sethp-org/sdk/loader"
)

type Client struct {
	serverID int
	token    string
	host     string
	http     *http.Client
}

func New(host string, token string, serverID int) *Client {
	return &Client{
		serverID: serverID,
		token:    token,
		host:     host,
		http: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type LogResult struct {
	Logs        []loader.Log    `json:"logs"`
	MetadataRaw json.RawMessage `json:"metadata"`
}

type LoadParams struct {
	Nickname string      `json:"nickname"`
	Date     time.Time   `json:"date"`
	Dates    []time.Time `json:"dates"`
}

func (c *Client) loadLogs(element string, param LoadParams) (*LogResult, error) {
	var paramMap = make(map[string]interface{})
	if param.Nickname != "" {
		paramMap["nickname"] = param.Nickname
	}
	if !param.Date.IsZero() {
		paramMap["date"] = param.Date
	}
	if len(param.Dates) == 2 {
		paramMap["dates"] = param.Dates
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(paramMap); err != nil {
		return nil, err
	}

	var result LogResult
	req, err := http.NewRequest(http.MethodPost, c.host+fmt.Sprintf("/logs/%s/%d", element, c.serverID), &buf)
	if err != nil {
		return nil, err
	}
	fmt.Println(req.URL.String())
	defer req.Body.Close()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type LogReportMetadata struct {
	Count int `json:"count"`
}

type LogReportResult struct {
	Logs     []loader.Log      `json:"logs"`
	Metadata LogReportMetadata `json:"metadata"`
}

type LogOnlineMetadata struct {
	Count        int    `json:"count"`
	LogoutCount  int    `json:"logout_count"`
	Online       int    `json:"online"`
	OnlineFormat string `json:"online_format"`
}

type LogOnlineResult struct {
	Logs     []loader.Log      `json:"logs"`
	Metadata LogOnlineMetadata `json:"metadata"`
}

func (c *Client) Report(param LoadParams) (*LogReportResult, error) {
	result, err := c.loadLogs("report", param)
	if err != nil {
		return nil, err
	}
	var metadata LogReportMetadata
	if err := json.Unmarshal(result.MetadataRaw, &metadata); err != nil {
		return nil, err
	}
	result.MetadataRaw = nil
	return &LogReportResult{
		Logs:     result.Logs,
		Metadata: metadata,
	}, nil
}

func (c *Client) Online(param LoadParams) (*LogOnlineResult, error) {
	result, err := c.loadLogs("online", param)
	if err != nil {
		return nil, err
	}
	var metadata LogOnlineMetadata
	if err := json.Unmarshal(result.MetadataRaw, &metadata); err != nil {
		return nil, err
	}
	result.MetadataRaw = nil
	return &LogOnlineResult{
		Logs:     result.Logs,
		Metadata: metadata,
	}, nil
}

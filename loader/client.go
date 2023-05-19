package loader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	HostAz = "az"
	HostRD = "rd"
)

type Client struct {
	settings Settings
	client   *http.Client
}

type Settings struct {
	Host     string
	BaseHost string
	Token    string
	Login    string
	Password string
	Secret   string
	ServerID int
	Timeout  time.Duration
}

func New(settings Settings) *Client {
	return &Client{
		settings: settings,
		client: &http.Client{
			Timeout: settings.Timeout,
		},
	}
}

func (c *Client) addUlogParams(ulogParameters UlogParameters) (parameters map[string]interface{}) {
	parameters = make(map[string]interface{})
	parameters["user_id"] = ulogParameters.UserID
	parameters["is_vc"] = ulogParameters.IsVC
	parameters["is_admin"] = ulogParameters.IsAdmin
	parameters["text"] = ulogParameters.Text
	parameters["full"] = ulogParameters.Full
	parameters["date"] = ulogParameters.Date
	parameters["listen"] = ulogParameters.Listen
	parameters["type"] = ulogParameters.Type
	return
}

func (c *Client) addGamePanelParams(gamePanelParameters GamePanelParameters) (parameters map[string]interface{}) {
	parameters = make(map[string]interface{})
	parameters["host"] = c.settings.BaseHost
	parameters["secret"] = c.settings.Secret
	parameters["nickname"] = gamePanelParameters.Nickname
	parameters["params"] = gamePanelParameters.Params
	parameters["reason"] = gamePanelParameters.Reason
	if gamePanelParameters.Dates != nil {
		parameters["dates"] = gamePanelParameters.Dates
	}
	if !gamePanelParameters.Date.IsZero() {
		parameters["date"] = gamePanelParameters.Date
	}
	parameters["listen"] = gamePanelParameters.Listen
	parameters["type"] = gamePanelParameters.Type
	return parameters
}

func (c *Client) doRequest(engine string, headers, parameters map[string]interface{}) (io.ReadCloser, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(parameters); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.settings.Host+"/api/"+engine+"/parse", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		var bufErr bytes.Buffer
		if _, err := io.Copy(&bufErr, resp.Body); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("download error: %s", bufErr.String())
	}

	return resp.Body, nil
}

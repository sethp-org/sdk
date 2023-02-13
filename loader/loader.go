package loader

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const (
	FormatRaw  = "raw"
	FormatJson = "json"
)

const (
	EngineUlog      = "ulog"
	EngineGamePanel = "gamepanel"
)

var ErrUnknownEngine = errors.New("unknown engine")

func (c *Client) downloadRaw(r io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Client) downloadLogJSON(r io.ReadCloser) ([]Log, error) {
	var logs []Log
	if err := json.NewDecoder(r).Decode(&logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (c *Client) downloadUserJSON(r io.ReadCloser) ([]User, error) {
	var users []User
	if err := json.NewDecoder(r).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *Client) download(format string, parameters interface{}) (io.ReadCloser, error) {
	var (
		body io.ReadCloser
		err  error
	)

	var params = make(map[string]interface{})

	var headers = map[string]interface{}{
		"token":     c.settings.Token,
		"login":     c.settings.Login,
		"password":  c.settings.Password,
		"base_host": c.settings.BaseHost,
		"secret":    c.settings.Secret,
	}
	var engine string

	switch p := parameters.(type) {
	case UlogParameters:
		params = c.addUlogParams(p)
		engine = EngineUlog
	case GamePanelParameters:
		params = c.addGamePanelParams(p)
		engine = EngineGamePanel
	default:
		return nil, ErrUnknownEngine
	}

	params["format"] = format
	params["server_id"] = c.settings.ServerID

	body, err = c.doRequest(engine, headers, params)
	if err != nil {
		return nil, fmt.Errorf("doRequest: %w", err)
	}
	return body, nil
}

func (c *Client) DownloadRaw(parameters interface{}) (string, error) {
	body, err := c.download(FormatRaw, parameters)
	if err != nil {
		return "", err
	}
	defer body.Close()
	return c.downloadRaw(body)
}

func (c *Client) DownloadLogJSON(parameters interface{}) ([]Log, error) {
	body, err := c.download(FormatJson, parameters)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return c.downloadLogJSON(body)
}

func (c *Client) DownloadUserJSON(parameters interface{}) ([]User, error) {
	body, err := c.download(FormatJson, parameters)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return c.downloadUserJSON(body)
}

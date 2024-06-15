package oauth2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DiscordAuthUrl   = "https://discord.com/oauth2/authorize"
	DiscordTokenURL  = "https://discord.com/api/oauth2/token"
	DiscordUserMeURL = "https://discord.com/api/users/@me"
)

type DiscordConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scope        string
	Timeout      time.Duration
}

type DiscordOAuth2 struct {
	clientID     string
	clientSecret string
	redirectURL  string
	scope        string

	client *http.Client
}

type DiscordUserInfo struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
}

func NewDiscord(cfg DiscordConfig) *DiscordOAuth2 {
	return &DiscordOAuth2{
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		redirectURL:  cfg.RedirectURL,
		scope:        cfg.Scope,

		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (o *DiscordOAuth2) AuthURL() string {
	u := url.Values{}
	u.Set("client_id", o.clientID)
	u.Set("response_type", "code")
	u.Set("scope", o.scope)
	u.Set("redirect_uri", o.redirectURL)
	return DiscordAuthUrl + "?" + u.Encode()
}

func (o *DiscordOAuth2) Token(code string) (*Token, error) {
	u := url.Values{}
	u.Set("client_id", o.clientID)
	u.Set("client_secret", o.clientSecret)
	u.Set("grant_type", "authorization_code")
	u.Set("code", code)
	u.Set("redirect_uri", o.redirectURL)

	var body io.Reader = strings.NewReader(u.Encode())
	req, err := http.NewRequest(http.MethodPost, DiscordTokenURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (o *DiscordOAuth2) UserInfo(token string) (*DiscordUserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, DiscordUserMeURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo DiscordUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	if userInfo.ID == "" {
		return nil, fmt.Errorf("user not found")
	}

	return &userInfo, nil
}

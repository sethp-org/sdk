package oauth2

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	WellKnowURL  string
	Scope        string
	Timeout      time.Duration
}

type OAuth2 struct {
	clientID     string
	clientSecret string
	redirectURL  string
	wellKnowURL  string
	scope        string

	client *http.Client
}

func New(cfg Config) *OAuth2 {
	return &OAuth2{
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		redirectURL:  cfg.RedirectURL,
		wellKnowURL:  cfg.WellKnowURL,
		scope:        cfg.Scope,

		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

type WellKnow struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserInfoEndpoint      string `json:"userinfo_endpoint"`
}

func (o *OAuth2) WellKnow() (WellKnow, error) {
	resp, err := o.client.Get(o.wellKnowURL)
	if err != nil {
		return WellKnow{}, err
	}
	defer resp.Body.Close()

	var data WellKnow
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return WellKnow{}, err
	}

	return data, nil
}

func (o *OAuth2) AuthURL() (string, error) {
	data, err := o.WellKnow()
	if err != nil {
		return "", err
	}
	u := url.Values{}
	u.Set("client_id", o.clientID)
	u.Set("response_type", "code")
	u.Set("scope", o.scope)
	u.Set("redirect_uri", o.redirectURL)
	return data.AuthorizationEndpoint + "?" + u.Encode(), nil
}

type Token struct {
	AccessToken string `json:"access_token"`
}

func (o *OAuth2) Token(code string) (*Token, error) {
	data, err := o.WellKnow()
	if err != nil {
		return nil, err
	}

	u := url.Values{}
	u.Set("client_id", o.clientID)
	u.Set("client_secret", o.clientSecret)
	u.Set("grant_type", "authorization_code")
	u.Set("code", code)
	u.Set("redirect_uri", o.redirectURL)

	var body io.Reader = strings.NewReader(u.Encode())
	req, err := http.NewRequest(http.MethodPost, data.TokenEndpoint, body)
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

type UserInfo struct {
	Sub               string `json:"sub"`
	EmailVerified     bool   `json:"email_verified"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Picture           string `json:"picture"`
	PreferredUsername string `json:"preferred_username"`
}

func (o *OAuth2) UserInfo(token string) (*UserInfo, error) {
	data, err := o.WellKnow()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, data.UserInfoEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

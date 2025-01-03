package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Provider interface {
	GetAuthURL() string
	GetUserInfo(code string) (*UserInfo, error)
}

type UserInfo struct {
	ID       string
	Email    string
	Name     string
	Provider string
}

type GoogleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider(clientID, clientSecret, redirectURL string) *GoogleProvider {
	return &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (p *GoogleProvider) GetAuthURL() string {
	return p.config.AuthCodeURL("state")
}

func (p *GoogleProvider) GetUserInfo(code string) (*UserInfo, error) {
	token, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %v", err)
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed parsing response: %v", err)
	}

	return &UserInfo{
		ID:       result["id"].(string),
		Email:    result["email"].(string),
		Name:     result["name"].(string),
		Provider: "google",
	}, nil
}

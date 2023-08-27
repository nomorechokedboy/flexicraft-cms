package auth

import (
	"context"
	"encoding/json"
	"strconv"

	"golang.org/x/oauth2"
)

const GitlabName = "gitlab"

type Gitlab struct {
	*baseProvider
}

func NewGitlabProvider() *Gitlab {
	return &Gitlab{&baseProvider{
		authUrl:      "https://gitlab.com/oauth/authorize",
		clientId:     "a88c7075c5f5c029d5803f0f6d08490140e9d1deb318b270db04194ccc3b1527",
		clientSecret: "",
		ctx:          context.Background(),
		redirectUrl:  "http://localhost:3000/gitlab",
		scopes:       []string{"read_user"},
		userApiUrl:   "https://gitlab.com/api/v4/user",
		tokenUrl:     "https://gitlab.com/oauth/token",
	}}
}

func (p *Gitlab) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawAuthUser(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatar_url"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           strconv.Itoa(extracted.Id),
		Name:         extracted.Name,
		Username:     extracted.Username,
		Email:        extracted.Email,
		AvatarUrl:    extracted.AvatarUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return user, nil
}

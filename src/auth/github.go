package auth

import (
	"context"
	"encoding/json"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var _ Provider = (*Github)(nil)

const GithubName = "github"

type Github struct {
	*baseProvider
}

func NewGithubProvider() *Github {
	return &Github{
		&baseProvider{
			authUrl:      github.Endpoint.AuthURL,
			clientId:     "a7212b06f4ba590d04ee",
			clientSecret: "",
			ctx:          context.Background(),
			redirectUrl:  "http://localhost:3000/callback",
			scopes:       []string{"read:user", "user:mail"},
			userApiUrl:   "https://api.github.com/user",
			tokenUrl:     github.Endpoint.TokenURL,
		},
	}
}

func (p *Github) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawAuthUser(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Login     string `json:"login"`
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatar_url"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           strconv.Itoa(extracted.Id),
		Name:         extracted.Name,
		Username:     extracted.Login,
		Email:        extracted.Email,
		AvatarUrl:    extracted.AvatarUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return user, nil
}

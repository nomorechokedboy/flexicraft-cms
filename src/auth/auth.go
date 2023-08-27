package auth

import (
	"fmt"

	"golang.org/x/oauth2"
)

type AuthUser struct {
	AccessToken  string         `json:"accessToken"`
	AvatarUrl    string         `json:"avatarUrl"`
	Email        string         `json:"email"`
	Id           string         `json:"id"`
	Name         string         `json:"name"`
	RawUser      map[string]any `json:"rawUser"`
	RefreshToken string         `json:"refreshToken"`
	Username     string         `json:"username"`
}

type Provider interface {
	FetchToken(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	FetchRawAuthUser(token *oauth2.Token) ([]byte, error)
	FetchAuthUser(token *oauth2.Token) (*AuthUser, error)
}

func NewProviderByName(name string) (Provider, error) {
	switch name {
	case GithubName:
		return NewGithubProvider(), nil
	case GitlabName:
		return NewGitlabProvider(), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

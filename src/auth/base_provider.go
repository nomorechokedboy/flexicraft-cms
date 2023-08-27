package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type baseProvider struct {
	authUrl      string
	clientId     string
	clientSecret string
	ctx          context.Context
	redirectUrl  string
	scopes       []string
	tokenUrl     string
	userApiUrl   string
}

func (bp *baseProvider) oauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     bp.clientId,
		ClientSecret: bp.clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  bp.authUrl,
			TokenURL: bp.tokenUrl,
		},
		RedirectURL: bp.redirectUrl,
		Scopes:      bp.scopes,
	}
}

func (bp *baseProvider) FetchToken(
	code string,
	opts ...oauth2.AuthCodeOption,
) (*oauth2.Token, error) {
	fmt.Printf("Code: %s, options: %#v\n", code, opts)
	return bp.oauth2Config().Exchange(bp.ctx, code, opts...)
}

func (bp *baseProvider) sendAuthUserDataRequest(
	req *http.Request,
	token *oauth2.Token,
) ([]byte, error) {
	client := http.Client{Timeout: time.Second * 30}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	defer res.Body.Close()
	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (bp *baseProvider) FetchRawAuthUser(token *oauth2.Token) ([]byte, error) {
	req, err := http.NewRequestWithContext(bp.ctx, "GET", bp.userApiUrl, nil)
	if err != nil {
		return nil, err
	}

	return bp.sendAuthUserDataRequest(req, token)
}

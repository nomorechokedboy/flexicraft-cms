package usecases

import (
	apperr "api/src/app_error"
	"api/src/auth"
	"api/src/entities"
	"api/src/tokenizer"
	"fmt"
	"log"
	"time"
)

type (
	OAuth interface {
		Execute(code string, provider string) (*entities.AuthResponse, error)
	}

	BaseOAuth struct {
		tokenizer tokenizer.Tokenizer
	}
)

var _ OAuth = (*BaseOAuth)(nil)

func NewOAuth(tokenizer tokenizer.Tokenizer) OAuth {
	return &BaseOAuth{tokenizer}
}

func (o *BaseOAuth) Execute(code string, provider string) (*entities.AuthResponse, error) {
	oauthProvider, err := auth.NewProviderByName(provider)
	fmt.Printf("Provider: %#v\n", oauthProvider)
	if err != nil {
		return nil, apperr.New("100005", 400, "Unknown provider", "Bad request", err)
	}

	token, err := oauthProvider.FetchToken(code)
	if err != nil {
		log.Println("Failed to get token:", err)
		return nil, apperr.New(
			"100011",
			500,
			"Failed to get oauth token",
			"Internal server error",
			err,
		)
	}

	user, err := oauthProvider.FetchAuthUser(token)
	if err != nil {
		log.Println("Failed to get user:", err)
		return nil, apperr.New(
			"1000012",
			500,
			"Failed to get user information",
			"Internal server error",
			err,
		)
	}
	log.Printf("user: %#v\n", user)

	tokenPayload := tokenizer.Payload{Identifier: user.Username}
	accessToken, err := o.tokenizer.Sign(
		tokenPayload,
		[]byte("token-secret"),
		time.Duration(1000*60*60*5),
	)
	if err != nil {
		return nil, err
	}
	refreshToken, err := o.tokenizer.Sign(
		tokenPayload,
		[]byte("refresh-token-secret"),
		time.Duration(1000*60*60*24*30),
	)
	if err != nil {
		return nil, err
	}

	return &entities.AuthResponse{Token: *accessToken, RefreshToken: *refreshToken}, nil
}

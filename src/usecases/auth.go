package usecases

import (
	apperr "api/src/app_error"
	"api/src/entities"
	"api/src/hasher"
	"api/src/repositories"
	"api/src/tokenizer"
	"api/src/validator"
	"time"
)

type (
	Authenticator interface {
		SignUp(entities.CreateAuth) (*entities.AuthEntity, error)
		SignIn(entities.CreateAuth) (*entities.AuthResponse, error)
	}

	BaseAuthenticator struct {
		hasher    hasher.Hasher
		repo      repositories.AuthRepo
		tokenizer tokenizer.Tokenizer
		validator validator.Validator
	}
)

var _ Authenticator = (*BaseAuthenticator)(nil)

func New(
	hasher hasher.Hasher,
	repo repositories.AuthRepo,
	tokenizer tokenizer.Tokenizer,
	validator validator.Validator,
) Authenticator {
	return &BaseAuthenticator{hasher, repo, tokenizer, validator}
}

func (ba *BaseAuthenticator) SignUp(payload entities.CreateAuth) (*entities.AuthEntity, error) {
	if err := ba.validator.ValidateStruct(&payload); err != nil {
		return nil, apperr.New("100000", 400, "Validate failed", "Bad request", err)
	}

	hashedPassword, err := ba.hasher.Hash(payload.Password)
	if err != nil {
		return nil, apperr.New("100001", 500, "Hash error", "Internal server error", err)
	}

	payload.Password = *hashedPassword
	res, err := ba.repo.Insert(payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ba *BaseAuthenticator) SignIn(payload entities.CreateAuth) (*entities.AuthResponse, error) {
	authEntity := entities.AuthEntity{Identifier: payload.Identifier}
	res, err := ba.repo.FindOne(authEntity)
	if err != nil {
		return nil, err
	}

	match, err := ba.hasher.Verify(payload.Password, res.Password)
	if err != nil {
		return nil, err
	}

	if !*match {
		return nil, apperr.New("102001", 400, "Wrong username or password", "Bad request", nil)
	}

	tokenPayload := tokenizer.Payload{Identifier: res.Identifier}
	token, err := ba.tokenizer.Sign(
		tokenPayload,
		[]byte("token-secret"),
		time.Duration(1000*60*60*5),
	)
	if err != nil {
		return nil, err
	}
	refreshToken, err := ba.tokenizer.Sign(
		tokenPayload,
		[]byte("refresh-token-secret"),
		time.Duration(1000*60*60*24*30),
	)
	if err != nil {
		return nil, err
	}

	return &entities.AuthResponse{Token: *token, RefreshToken: *refreshToken}, nil
}

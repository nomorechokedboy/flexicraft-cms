package tokenizer

import (
	apperr "api/src/app_error"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	Tokenizer interface {
		Sign(Payload, []byte, time.Duration) (*string, error)
		Verify(string, []byte) (*Claims, error)
	}

	Payload struct {
		Identifier string
	}

	Claims struct {
		Payload
		jwt.RegisteredClaims
	}

	Jwt struct{}
)

var _ Tokenizer = (*Jwt)(nil)

func (j Jwt) Sign(payload Payload, secret []byte, duration time.Duration) (*string, error) {
	t := time.Now()
	expiresAt := t.Add(duration * time.Millisecond)
	claims := &Claims{
		Payload:          payload,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expiresAt)},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, apperr.New("100005", 500, "Can't sign token", "Internal server error", err)
	}

	return &tokenString, nil
}

func (j Jwt) Verify(tokenString string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, apperr.New("100004", 401, "Invalid signature", "Unauthorized", err)
		}

		return nil, apperr.New("100002", 400, "Bad request", "Bad request", err)
	}
	if !token.Valid {
		return nil, apperr.New("100004", 401, "Invalid signature", "Unauthorized", err)
	}

	return claims, nil
}

package hasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type (
	Argon2HashParams struct {
		Memory      uint32
		Iterations  uint32
		Parallelism uint8
		SaltLength  uint32
		KeyLength   uint32
	}

	Argon2Hasher struct {
		Argon2HashParams
	}
)

var _ Hasher = (*Argon2Hasher)(nil)

func (a2 Argon2Hasher) Hash(password string) (*string, error) {
	salt, err := a2.genSalt(a2.SaltLength)
	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		a2.Iterations,
		a2.Memory,
		a2.Parallelism,
		a2.KeyLength,
	)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		a2.Memory,
		a2.Iterations,
		a2.Parallelism,
		b64Salt,
		b64Hash,
	)

	return &encodedHash, nil
}

func (a2 Argon2Hasher) Verify(password string, encodedHash string) (*bool, error) {
	params, salt, hash, err := a2.decodeHash(encodedHash)
	if err != nil {
		return nil, err
	}

	res := false
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		res = true
	}

	return &res, nil
}

func (a2 *Argon2Hasher) genSalt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func (a2 *Argon2Hasher) decodeHash(decodedHash string) (*Argon2HashParams, []byte, []byte, error) {
	vals := strings.Split(decodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	if _, err := fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	params := &Argon2HashParams{}
	if _, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism); err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(hash))

	return params, salt, hash, nil
}

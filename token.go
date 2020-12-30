package vatinator

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto/v2"
)

type TokenService interface {
	NewPath(path string) (string, error)
	CheckPath(encToken string, path string) error
}

type tokenService struct {
	key []byte
}

func NewTokenService(key []byte) TokenService {
	return tokenService{key}
}

func (ts tokenService) NewPath(path string) (string, error) {
	token := paseto.JSONToken{
		Subject:    path,
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(30 * 24 * time.Hour),
	}
	encToken, err := paseto.Encrypt(ts.key, token, "")
	if err != nil {
		return "", err
	}
	return encToken, nil
}

func (ts tokenService) CheckPath(encToken string, path string) error {
	var token paseto.JSONToken
	var footer string
	if err := paseto.Decrypt(encToken, ts.key, &token, &footer); err != nil {
		return err
	}
	if token.Subject != path {
		return fmt.Errorf("not valid for path %s: got %s", path, token.Subject)
	}
	return nil
}

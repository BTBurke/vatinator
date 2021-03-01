package vatinator

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto/v2"
)

const passwordSubject = "password-reset"

type TokenService interface {
	NewPath(path string) (string, error)
	CheckPath(encToken string, path string) error
	NewPasswordReset(email string) (string, error)
	CheckPasswordReset(encToken string) (string, error)
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

func (ts tokenService) NewPasswordReset(email string) (string, error) {
	token := paseto.JSONToken{
		Subject:    passwordSubject,
		Audience:   email,
		Expiration: time.Now().Add(30 * time.Minute),
		IssuedAt:   time.Now(),
	}
	encToken, err := paseto.Encrypt(ts.key, token, "")
	if err != nil {
		return "", err
	}
	return encToken, nil
}

func (ts tokenService) CheckPasswordReset(encToken string) (string, error) {
	var token paseto.JSONToken
	var footer string
	if err := paseto.Decrypt(encToken, ts.key, &token, &footer); err != nil {
		return "", err
	}
	if token.Subject != passwordSubject || len(token.Audience) == 0 {
		return "", fmt.Errorf("not a valid password reset token, subject=%s, len(email)=%d", token.Subject, len(token.Audience))
	}
	return token.Audience, nil
}

var _ TokenService = tokenService{}

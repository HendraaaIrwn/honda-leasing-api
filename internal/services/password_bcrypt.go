package services

import (
	"strings"

	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

func hashPasswordWithBcrypt(password string) (string, error) {
	password = strings.TrimSpace(password)
	if password == "" {
		return "", errs.ErrInvalidPassword
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func comparePasswordWithBcrypt(hashedPassword, password string) error {
	hashedPassword = strings.TrimSpace(hashedPassword)
	password = strings.TrimSpace(password)

	if hashedPassword == "" || password == "" {
		return errs.ErrInvalidPassword
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

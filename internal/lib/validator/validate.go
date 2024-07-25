package validator

import (
	"fmt"
	"strings"
)

// ValidatePassword проверяет корректность пароля указанным требованиям.
// Если все ок, то возвращает nil, иначе - ошибку.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password cannot be shorter than 8 characters")
	}

	if fields := strings.Fields(password); len(fields) > 1 {
		return fmt.Errorf("password cannot contain spaces")
	}

	return nil
}

// ValidateLogin проверяет корректность логина указанным требованиям.
// Если все ок, то возвращает nil, иначе - ошибку.
func ValidateLogin(login string) error {
	if len(login) < 8 {
		return fmt.Errorf("login cannot be shorter than 8 characters")
	}

	if fields := strings.Fields(login); len(fields) > 1 {
		return fmt.Errorf("login cannot contain spaces")
	}

	return nil
}

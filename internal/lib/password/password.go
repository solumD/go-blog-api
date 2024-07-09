package password

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	cost := 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashAndPass(password, realPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(realPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

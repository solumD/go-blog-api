package password

import "golang.org/x/crypto/bcrypt"

// EncryptPassword возвращает зашифрованый пароль
func EncryptPassword(password string) (string, error) {
	cost := 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CompareHashAndPass сравнивает зашифрованный пароль с другим паролем
// в случае несовпадения возвращает ошибку, иначе - nil
func CompareHashAndPass(password, realPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(realPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

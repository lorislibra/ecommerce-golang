package utils

import "golang.org/x/crypto/bcrypt"

func GenPasswordHash(password string) (string, error) {
	blob, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(blob), err
}

func TestPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

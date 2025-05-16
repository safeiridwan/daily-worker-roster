package user

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"math/big"
)

const (
	passwordLength = 8
	charset        = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"
)

func generateRandomPassword() (string, error) {
	password := make([]byte, passwordLength)
	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[n.Int64()]
	}
	return string(password), nil
}

func GeneratePassword() (plainPassword string, hashedPassword string, err error) {
	plainPassword, err = generateRandomPassword()
	if err != nil {
		return "", "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	hashedPassword = string(hashedBytes)
	return plainPassword, hashedPassword, nil
}

package passwordgenerator

import (
	"crypto/rand"
	"math/big"
)

func generateRandomNumber(max int64) (int64, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, err
	}
	return num.Int64(), nil
}

// Function to generate a secure random alphanumeric password
func GenerateSecureRandomPassword(length int, numbers, symbols, uppercase bool) (string, error) {
	var passwordSet string
	if numbers {
		passwordSet += "0123456789"
	}
	if symbols {
		passwordSet += "!@#$%^&*()"
	}
	if uppercase {
		passwordSet += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	passwordSet += "abcdefghijklmnopqrstuvwxyz"

	password := make([]byte, length)
	for i := range password {
		num, err := generateRandomNumber(int64(len(passwordSet)))
		if err != nil {
			return "", err
		}
		password[i] = passwordSet[num]
	}

	return string(password), nil
}

// Function to generate a secure random alphanumeric password
func GenerateSecureAlphanumericPassword(length int, numbers, uppercase bool) (string, error) {
	var passwordSet string
	if numbers {
		passwordSet += "0123456789"
	}
	if uppercase {
		passwordSet += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	passwordSet += "abcdefghijklmnopqrstuvwxyz"

	password := make([]byte, length)
	for i := range password {
		num, err := generateRandomNumber(int64(len(passwordSet)))
		if err != nil {
			return "", err
		}
		password[i] = passwordSet[num]
	}

	return string(password), nil
}

// Function to generate a secure random PIN
func GenerateSecurePIN(length int) (string, error) {
	pin := make([]byte, length)
	for i := range pin {
		num, err := generateRandomNumber(10)
		if err != nil {
			return "", err
		}
		pin[i] = byte(num + 48) // ASCII digits start from 48
	}

	return string(pin), nil
}

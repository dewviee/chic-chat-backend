package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RandomStringFromText(input string) (string, error) {
	// Convert the input string to a byte slice
	data := []byte(input)

	// Generate a random byte slice of the same length as the input
	randomBytes := make([]byte, len(data))
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// XOR the input with the random bytes to create the random string
	for i := 0; i < len(data); i++ {
		randomBytes[i] ^= data[i]
	}

	// Encode the random bytes to base64 to get the final random string
	randomString := base64.URLEncoding.EncodeToString(randomBytes)
	return randomString, nil
}

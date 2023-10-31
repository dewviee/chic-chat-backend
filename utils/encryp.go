package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
)

// ฟังก์ชันเข้ารหัส RSA
func EncryptRSA(publicKey *rsa.PublicKey, plaintext string) (string, error) {
	plaintextBytes := []byte(plaintext)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintextBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(ciphertext), nil
}

// ฟังก์ชันถอดรหัส RSA
func DecryptRSA(privateKey *rsa.PrivateKey, ciphertext string) (string, error) {
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	decrypted, err := privateKey.Decrypt(rand.Reader, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

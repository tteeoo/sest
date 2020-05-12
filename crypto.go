package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"io"
)

func a2Hash(password string, salt []byte) []byte {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return hash
}

func generateSalt(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func encrypt(data string, key []byte) ([]byte, error) {
	cphr, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, []byte(data), nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcmDecrypt, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcmDecrypt.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, err
	}

	nonce, encryptedMessage := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcmDecrypt.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func bEncode(b []byte) string {
	return base64.RawStdEncoding.EncodeToString(b)
}

func bDecode(s string) ([]byte, error) {
	b, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

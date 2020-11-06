package utils

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/scrypt"
)

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

func GenerateSalt() (string, error) {
	salt := make([]byte, PW_SALT_BYTES)
	_, err := io.ReadFull(rand.Reader, salt)

	return string(salt), err
}

func HashPassword(password, salt string) (string, error) {
	hash, err := scrypt.Key([]byte(password), []byte(salt), 1<<14, 8, 1, PW_HASH_BYTES)

	return string(hash), err
}

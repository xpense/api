package utils

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/scrypt"
)

type PasswordHasher interface {
	GenerateSalt() (string, error)
	HashPassword(password, salt string) (string, error)
}

const (
	pwSaltBytes = 32
	pwHashBytes = 64
)

type hasher struct{}

func NewPasswordHasher() PasswordHasher {
	return &hasher{}
}

// GenerateSalt generates a salt for hashing passwords
func (h *hasher) GenerateSalt() (string, error) {
	salt := make([]byte, pwSaltBytes)
	_, err := io.ReadFull(rand.Reader, salt)

	return hex.EncodeToString(salt), err
}

// HashPassword hashes a given password with a salt
func (h *hasher) HashPassword(password, salt string) (string, error) {
	saltSlice, err := hex.DecodeString(salt)
	if err != nil {
		return "", err
	}

	hash, err := scrypt.Key([]byte(password), saltSlice, 1<<14, 8, 1, pwHashBytes)
	return hex.EncodeToString(hash), err
}

package utils

import (
	"bytes"
	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

const (
	saltLen = 8

	time    = 1
	memory  = 64 * 1024
	threads = 4
	keyLen  = 32
)

func HashPass(plainPassword string) ([]byte, error) {
	salt := make([]byte, saltLen)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return hashPassWithSalt(salt, plainPassword), nil
}

func hashPassWithSalt(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, time, memory, threads, keyLen)

	return append(salt, hashedPass...)
}

func ComparePassAndHash(passHash []byte, plainPassword string) bool {
	salt := passHash[0:8]
	userPassHash := hashPassWithSalt(salt, plainPassword)

	return bytes.Equal(userPassHash, passHash)
}

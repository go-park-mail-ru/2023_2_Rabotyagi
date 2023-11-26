package usecases

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashContent(content []byte) (string, error) {
	hash := sha256.New()

	_, err := hash.Write(content)
	if err != nil {
		return "", err
	}

	result := hash.Sum(nil)

	return hex.EncodeToString(result), nil
}

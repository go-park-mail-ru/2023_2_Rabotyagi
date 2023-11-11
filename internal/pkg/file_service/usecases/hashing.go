package usecases

import (
	"crypto/sha256"
)

func HashContent(content []byte) string {
	hash := sha256.New()
	result := hash.Sum(content)

	return string(result)
}

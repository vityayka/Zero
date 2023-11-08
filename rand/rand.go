package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const sessionTokenSize = 32

func bytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func String(size int) (string, error) {
	randBytes, err := bytes(size)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randBytes), nil
}

func SessionToken() (string, error) {
	return String(sessionTokenSize)
}

func IsSessionToken(input string) bool {
	return len(input) == sessionTokenSize
}

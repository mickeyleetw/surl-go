package utils

import (
	"crypto/md5"
	"encoding/base64"
	"strings"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = 6
)

func GenerateShortURL(longURL string) string {

	hasher := md5.New()
	hasher.Write([]byte(longURL))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	shortURL := strings.Builder{}
	for i := 0; i < length; i++ {
		c := hash[i]
		if strings.ContainsRune(alphabet, rune(c)) {
			shortURL.WriteByte(c)
		} else {
			shortURL.WriteByte(alphabet[int(c)%len(alphabet)])
		}
	}

	return shortURL.String()
}

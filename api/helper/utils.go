package helper

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"

	"github.com/lib/pq"
)

func ExtractUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return email
}

func GenerateState(n uint) string {
	b := make([]byte, n/2)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func IsUniqueConstraintError(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == "23505"
	}
	return false
}

func FormatURL(path *string) *string {
	if path == nil || *path == "" {
		return nil
	}

	url := *path
	if strings.HasPrefix(url, "https://") {
		return &url
	}

	formatted := os.Getenv("BE_DOMAIN") + "/" + url
	return &formatted
}

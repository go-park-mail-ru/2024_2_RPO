package encrypt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// GenerateSessionID создает безопасный идентификатор сессии длиной 64 символа (256 бит)
func GenerateSessionID() string {
	randomBytes := make([]byte, 32) // 256 бит
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		return ""
	}

	hash := sha256.New()
	hash.Write(randomBytes)

	sessionID := hex.EncodeToString(hash.Sum(nil))

	return sessionID
}

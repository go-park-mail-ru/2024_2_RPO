package encrypt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/satori/uuid"
	"golang.org/x/crypto/bcrypt"
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

// SaltAndHashPassword принимает строку пароля и возвращает хешированный пароль.
func SaltAndHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword проверяет, удовлетворяет ли данный пароль хешу
func CheckPassword(password string, hash string) bool {
	passwordBytes := []byte(password)
	hashBytes := []byte(hash)

	err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)

	return err == nil
}

// GenerateCSRFToken генерирует безопасный CSRF-токен
func GenerateCSRFToken() string {
	return uuid.NewV4().String()
}

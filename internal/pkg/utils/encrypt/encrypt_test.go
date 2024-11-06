package encrypt

import (
	"testing"

	"github.com/satori/uuid"
)

// Тест для GenerateSessionID
func TestGenerateSessionID(t *testing.T) {
	sessionID := GenerateSessionID()

	if len(sessionID) != 64 {
		t.Errorf("ожидалась длина 64, но получена %d", len(sessionID))
	}
}

// Тест для SaltAndHashPassword и CheckPassword
func TestSaltAndHashPasswordAndCheckPassword(t *testing.T) {
	password := "Test@123"

	// Тестирование хеширования пароля
	hashedPassword, err := SaltAndHashPassword(password)
	if err != nil {
		t.Fatalf("ошибка хеширования пароля: %v", err)
	}

	if hashedPassword == password {
		t.Errorf("хеш пароля не должен совпадать с исходным паролем")
	}

	// Проверка правильного пароля
	if !CheckPassword(password, hashedPassword) {
		t.Errorf("проверка пароля не удалась")
	}

	// Проверка неправильного пароля
	wrongPassword := "WrongPassword"
	if CheckPassword(wrongPassword, hashedPassword) {
		t.Errorf("неправильный пароль прошел проверку")
	}
}

// Тест для GenerateCSRFToken
func TestGenerateCSRFToken(t *testing.T) {
	token := GenerateCSRFToken()

	u, err := uuid.FromString(token)
	if err != nil || u.String() != token {
		t.Errorf("сгенерированный CSRF-токен не является корректным UUID: %v", err)
	}
}

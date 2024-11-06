package logging

import (
	"bytes"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestSetupLogger(t *testing.T) {
	// Создаем временный файл для проверки записи JSON-логов
	tmpFile, err := os.CreateTemp("", "test-log.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Настраиваем логгер
	SetupLogger(tmpFile)

	// Проверяем, что вывод в консоль настроен правильно
	consoleBuffer := &bytes.Buffer{}
	log.SetOutput(consoleBuffer)

	log.Info("Test log to console")
	if !bytes.Contains(consoleBuffer.Bytes(), []byte("Test log to console")) {
		t.Errorf("Лог не записан в консоль")
	}

	// Проверяем, что запись в файл работает
	log.Info("Test log to file")
	fileContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Не удалось прочитать содержимое временного файла: %v", err)
	}

	if !bytes.Contains(fileContent, []byte("Test log to file")) {
		t.Errorf("Лог не записан в файл")
	}
}

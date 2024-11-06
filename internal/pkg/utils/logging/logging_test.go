package logging

import (
	"bytes"
	"context"
	"fmt"
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

func TestLogFunc(t *testing.T) {
	ctx := context.WithValue(context.Background(), "requestID", uint64(123456))

	Warn(ctx, "This is a warning message")
	Info(ctx, "This is an info message")
	Error(ctx, "This is an error message")
	Debug(ctx, "This is a debug message")

	Warnf(ctx, "Warning with format: %s", "formatted warning")
	Infof(ctx, "Info with format: %s", "formatted info")
	Errorf(ctx, "Error with format: %s", "formatted error")
	Debugf(ctx, "Debug with format: %s", "formatted debug")

	fmt.Println("All functions tested.")
}

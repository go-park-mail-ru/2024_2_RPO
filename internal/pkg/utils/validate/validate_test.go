package validate_test

import (
	"context"
	"testing"

	"RPO_back/internal/pkg/utils/validate"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestValidate(t *testing.T) {
	// Создаем контекст
	ctx := context.Background()

	// Создаем экземпляр структуры с некорректными данными
	testData := TestStruct{
		Name:  "",
		Email: "invalid-email",
	}

	// Вызываем функцию валидации
	err := validate.Validate(ctx, testData)

	// Проверяем, что ошибка не равна nil
	assert.NotNil(t, err)

	validationErrors := err.(validator.ValidationErrors)

	// Проверяем, что есть 2 ошибки валидации
	assert.Len(t, validationErrors, 2)

	// Проверяем конкретные ошибки
	for _, err := range validationErrors {
		if err.Field() == "Name" {
			assert.Equal(t, "required", err.Tag())
		}
		if err.Field() == "Email" {
			assert.Equal(t, "email", err.Tag())
		}
	}

	// Тест с корректными данными
	validData := TestStruct{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	// Вызываем функцию валидации
	err = validate.Validate(ctx, validData)

	// Проверяем, что ошибки нет
	assert.Nil(t, err)
}

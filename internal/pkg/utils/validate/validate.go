package validate

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func Validate(v interface{}) error {
	validate := validator.New()

	// Проверяем, что аргумент — это структура, используя рефлексию
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct but got %s", val.Kind())
	}

	if err := validate.Struct(v); err != nil {
		// Форматируем ошибки валидации
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("Error: Field '%s' failed on the '%s' tag\n", err.Field(), err.Tag())
		}
		return err
	}

	return nil
}

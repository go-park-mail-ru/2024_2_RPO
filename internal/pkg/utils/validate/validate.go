package validate

import (
	"RPO_back/internal/pkg/utils/logging"
	"context"

	"github.com/go-playground/validator/v10"
)

func Validate(ctx context.Context, v interface{}) error {
	validate := validator.New()

	if err := validate.Struct(v); err != nil {
		// Форматируем ошибки валидации
		for _, err := range err.(validator.ValidationErrors) {
			logging.Warnf(ctx, "Error: Field '%s' failed on the '%s' tag\n", err.Field(), err.Tag())
		}
		return err
	}

	return nil
}

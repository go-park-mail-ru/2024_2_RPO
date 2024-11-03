package validate

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

func Validate(v interface{}) error {
	validate := validator.New()

	if err := validate.Struct(v); err != nil {
		// Форматируем ошибки валидации
		for _, err := range err.(validator.ValidationErrors) {
			log.Warnf("Error: Field '%s' failed on the '%s' tag\n", err.Field(), err.Tag())
		}
		return err
	}

	return nil
}

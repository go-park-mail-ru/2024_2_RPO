package validate

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/pkg/utils/logging"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Прекомпилированные регулярные выражения
var (
	upperCaseRegex = regexp.MustCompile(`[A-Z]+`)
	lowerCaseRegex = regexp.MustCompile(`[a-z]+`)
	digitRegex     = regexp.MustCompile(`[0-9]+`)
	usernameRegex  = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)
	specialCharSet = "!#&.,?/\\(){}[]\"'`;:|<>*^%~"
)

type Validatable interface {
	Validate() error
}

func Validate(ctx context.Context, v interface{}) error {
	validate := validator.New()

	if err := validate.Struct(v); err != nil {
		// Форматируем ошибки валидации
		for _, err := range err.(validator.ValidationErrors) {
			logging.Warnf(ctx, "Error: Field '%s' failed on the '%s' tag\n", err.Field(), err.Tag())
		}
		return err
	}

	if v, ok := v.(Validatable); ok && v != nil {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("Validate: %w", err)
		}
	}

	return nil
}

func CheckPassword(password string) error {
	if len(password) < 8 || len(password) > 50 {
		return fmt.Errorf("%w: password should be between 8 and 50 characters", errs.ErrValidation)
	}

	if !upperCaseRegex.MatchString(password) {
		return fmt.Errorf("%w: password should contain upper case latin character", errs.ErrValidation)
	}

	if !lowerCaseRegex.MatchString(password) {
		return fmt.Errorf("%w: password should contain lower case latin character", errs.ErrValidation)
	}

	if !digitRegex.MatchString(password) {
		return fmt.Errorf("%w: password should contain digit", errs.ErrValidation)
	}

	if !strings.ContainsAny(password, specialCharSet) {
		return fmt.Errorf("%w: password must contain at least one special character", errs.ErrValidation)
	}

	return nil
}

func CheckUserName(name string) error {
	if len(name) < 3 || len(name) > 30 {
		return fmt.Errorf("%w: name must be between 3 and 30 characters", errs.ErrValidation)
	}

	done := usernameRegex.Match([]byte(name))

	if !done {
		return fmt.Errorf("%w: name contains forbidden characters. Only allowed A-Z, a-z, 0-9, '_', '-', '.'", errs.ErrValidation)
	}

	return nil
}

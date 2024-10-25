package event

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"regexp"
)

func ValidateUUID() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		_, err := uuid.Parse(fl.Field().String())
		return err == nil
	}
}

func ValidateUsername() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		usernameToValidate := fl.Field().String()
		usernameExpression := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]$`)
		return usernameExpression.MatchString(usernameToValidate)
	}
}

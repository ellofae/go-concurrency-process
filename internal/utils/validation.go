package utils

import "github.com/go-playground/validator/v10"

func NewValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}

func ValidatorErrors(err error) map[string]string {
	fields := map[string]string{}

	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}

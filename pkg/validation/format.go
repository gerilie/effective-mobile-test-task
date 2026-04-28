package validation

import (
	"github.com/go-playground/validator/v10"
)

func formatErrorsByName(ve validator.ValidationErrors) Errors {
	errors := make(Errors)

	for _, err := range ve {
		message := getErrorMessageForTag(err)

		name := err.Field()
		errors[name] = message
	}

	return errors
}

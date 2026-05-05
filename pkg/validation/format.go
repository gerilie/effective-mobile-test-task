package validation

import (
	"github.com/go-playground/validator/v10"
)

// formatErrorsByName converts a slice of validator.ValidationErrors into an Errors map.
//
// It iterates over validation errors, generates a user-friendly message
// for each error using getErrorMessageForTag, and maps it by the field name.
//
// The resulting map uses field names as keys and corresponding error messages as values.
func formatErrorsByName(ve validator.ValidationErrors) Resp {
	errors := make(Resp)

	for _, fe := range ve {
		message := getErrorMessageForTag(fe)

		name := fe.Field()
		errors[name] = message
	}

	return errors
}

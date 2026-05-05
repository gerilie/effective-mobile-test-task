package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// getErrorMessageForTag returns a user-friendly error message for a validation error.
//
// It attempts to resolve the error message by delegating to specific
// handlers based on the validation tag (e.g., required, string, number, date).
// The first non-empty message returned by these handlers is used.
//
// If no specific handler matches, getErrorMessageForTag returns a default
// formatted message containing the validation tag and field name.
func getErrorMessageForTag(fe validator.FieldError) string {
	message := getErrorMessageForRequiredTag(fe)
	if message != "" {
		return message
	}

	message = getErrorMessageForStringTag(fe)
	if message != "" {
		return message
	}

	message = getErrorMessageForNumberTag(fe)
	if message != "" {
		return message
	}

	message = getErrorMessageForDatetimeTag(fe)
	if message != "" {
		return message
	}

	return fmt.Sprintf(
		"Validation failed for the '%s' rule on field '%s'",
		fe.Tag(),
		fe.Field(),
	)
}

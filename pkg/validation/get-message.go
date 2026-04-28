package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

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

	message = getErrorMessageForDateTag(fe)
	if message != "" {
		return message
	}

	return fmt.Sprintf(
		"Validation failed for the '%s' rule on field '%s'",
		fe.Tag(),
		fe.Field(),
	)
}

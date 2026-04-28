package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func getErrorMessageForNumberTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "min":
		return fmt.Sprintf(
			"The field must be at least %s characters long",
			fe.Param(),
		)
	case "max":
		return fmt.Sprintf(
			"The field must be at most %s characters long",
			fe.Param(),
		)
	case "gt":
		return fmt.Sprintf("The field must be greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf(
			"The field must be greater than or equal to %s",
			fe.Param(),
		)
	}

	return ""
}

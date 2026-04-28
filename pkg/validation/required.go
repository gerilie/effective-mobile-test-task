package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func getErrorMessageForRequiredTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "The field is required"
	case "required_with":
		return fmt.Sprintf("The field is required when '%s' is present", fe.Param())
	}

	return ""
}

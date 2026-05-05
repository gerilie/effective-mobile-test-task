package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// getErrorMessageForRequiredTag returns a user-friendly error message
// for required-related validation errors.
//
// It inspects the validation tag from the provided validator.FieldError
// and maps it to a corresponding message.
//
// Supported tags:
//   - "required":      the field must be present
//   - "required_with": the field must be present when the specified field is present
//
// If the validation tag is not recognized, getErrorMessageForRequiredTag returns an empty string.
func getErrorMessageForRequiredTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s present", ValidationPrefix)
	case "required_with":
		return fmt.Sprintf(
			"%s present when '%s' is present",
			ValidationPrefix,
			fe.Param(),
		)
	}

	return ""
}

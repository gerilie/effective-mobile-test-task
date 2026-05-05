package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// getErrorMessageForDatetimeTag returns a user-friendly error message for date-related validation errors.
//
// It inspects the validation tag from the provided validator.FieldError
// and maps it to a corresponding message.
//
// Supported tags:
//   - "datetime": the field must follow the "MM-YYYY" format.
//
// If the validation tag is not recognized, getErrorMessageForDatetimeTag returns an empty string.
func getErrorMessageForDatetimeTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "datetime":
		return fmt.Sprintf("%s in the format 'MM-YYYY'", ValidationPrefix)
	}

	return ""
}

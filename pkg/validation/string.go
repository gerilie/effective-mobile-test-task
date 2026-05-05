package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// getErrorMessageForStringTag returns a user-friendly error message
// for string-related validation errors.
//
// It inspects the validation tag from the provided validator.FieldError
// and maps it to a corresponding message.
//
// Supported tags:
//   - "email": the field must be a valid email address
//   - "uuid4": the field must be a valid UUID (version 4)
//
// If the validation tag is not recognized, getErrorMessageForStringTag returns an empty string.
func getErrorMessageForStringTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "email":
		return fmt.Sprintf("%s a valid email address", ValidationPrefix)
	case "uuid4":
		return fmt.Sprintf("%s a valid UUID (version 4)", ValidationPrefix)
	}

	return ""
}

package validation

import "github.com/go-playground/validator/v10"

func getErrorMessageForStringTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "email":
		return "The field must be a valid email address"
	case "uuid4":
		return "The field must be a valid UUID"
	}

	return ""
}

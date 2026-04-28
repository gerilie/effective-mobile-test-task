package validation

import "github.com/go-playground/validator/v10"

func getErrorMessageForDateTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "datetime":
		return "The field must be in the format MM-YYYY"
	}

	return ""
}

package validation

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

// WriteErrors writes validation errors to the HTTP response in JSON format.
//
// It converts the provided validator.ValidationErrors into a structured response
// grouped by field names and writes it with HTTP status 400 (Bad Request).
//
// Response format:
//
//	{
//	  "fields": {
//	    "field_name": "error message",
//	    "another_field": "error message"
//	  }
//	}
//
// If writing the response fails, WriteErrors logs the error using the logger
// extracted from the context.
func WriteErrors(ctx context.Context, w http.ResponseWriter, ve validator.ValidationErrors) {
	log := logger.FromContext(ctx)
	resp := formatErrorsByName(ve)

	if err := httputil.WriteJSON(ctx, w, http.StatusBadRequest, resp); err != nil {
		log.Error(ctx, "write response", zap.Error(err))
	}
}

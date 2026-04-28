package validation

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

func WriteErrors(ctx context.Context, w http.ResponseWriter, ve validator.ValidationErrors) {
	log := logger.FromContext(ctx)
	errors := formatErrorsByName(ve)

	resp := Resp{
		Fields: errors,
	}

	if err := httputil.WriteJSON(ctx, w, http.StatusBadRequest, resp); err != nil {
		log.Error(ctx, "write response", zap.Error(err))
	}
}

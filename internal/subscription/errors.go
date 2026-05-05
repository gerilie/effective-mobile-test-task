package subscription

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/validation"
	"go.uber.org/zap"
)

var errDateOrder = errors.New("date order")

// handleServiceErrors maps domain/service errors to HTTP responses.
func handleServiceErrors(ctx context.Context, w http.ResponseWriter, err error) {
	log := logger.FromContext(ctx)

	if errors.Is(err, errDateOrder) {
		log.Error(ctx, "validate subscription", zap.Error(err))

		resp := validation.Resp{
			"start_date": fmt.Sprintf("%s before end date", validation.ValidationPrefix),
		}

		if err := httputil.WriteJSON(ctx, w, http.StatusBadRequest, resp); err != nil {
			log.Error(ctx, "write response", zap.Error(err))
		}

		return
	}

	httputil.HandleDefaultErrors(ctx, w, err)
}

// handleValidationErrors processes validation errors and writes HTTP response.
func handleValidationErrors(ctx context.Context, w http.ResponseWriter, err error) {
	if err, ok := err.(validator.ValidationErrors); ok {
		validation.WriteErrors(ctx, w, err)
	} else {
		http.Error(w, "validation failed", http.StatusBadRequest)
	}
}

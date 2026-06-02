package subscription

import (
	"context"
	"errors"
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
		log.Error("validate subscription", zap.Error(err))

		resp := validation.Resp{
			"start_date": validation.ErrDateOrder.Error(),
		}

		if err := httputil.WriteJSON(w, http.StatusBadRequest, resp); err != nil {
			log.Error("write response", zap.Error(err))
		}

		return
	}

	httputil.HandleDefaultErrors(ctx, w, err)
}

// handleValidationErrors processes validation errors and writes HTTP response.
func handleValidationErrors(ctx context.Context, w http.ResponseWriter, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		validation.WriteErrors(ctx, w, ve)

		return
	}

	http.Error(w, httputil.InternalServerErrorMsg, http.StatusInternalServerError)
}

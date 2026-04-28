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

func handleServiceErrors(ctx context.Context, w http.ResponseWriter, err error) {
	log := logger.FromContext(ctx)

	if errors.Is(err, errDateOrder) {
		log.Error(ctx, "validate subscription", zap.Error(err))

		resp := validation.Resp{
			Fields: validation.Errors{
				"start_date": "start date must be before end date",
			},
		}

		if err := httputil.WriteJSON(ctx, w, http.StatusBadRequest, resp); err != nil {
			log.Error(ctx, "write response", zap.Error(err))
		}

		return
	}

	httputil.HandleDefaultErrors(ctx, w, err)
}

func handleValidationErrors(ctx context.Context, w http.ResponseWriter, err error) {
	if err, ok := err.(validator.ValidationErrors); ok {
		validation.WriteErrors(ctx, w, err)
	} else {
		http.Error(w, "validation failed", http.StatusBadRequest)
	}
}

package ping

import (
	"net/http"
	"time"

	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

// Ping is a simple health check endpoint that returns a timestamp.
//
//	@Summary	Ping service
//	@Tags		health
//	@ID			ping
//	@Produce	json
//	@Success	200	{object}	Resp	"Ping success"
//	@Failure	500	{string}	string	"Internal server error"
//	@Router		/ping [get].
func Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	resp := Resp{
		Timestamp: time.Now(),
	}

	if err := httputil.WriteJSON(w, http.StatusOK, resp); err != nil {
		log.Error("ping failed", zap.Error(err))

		return
	}

	log.Info("ping success")
}

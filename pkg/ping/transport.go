package ping

import (
	"net/http"
	"time"

	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

// @Summary	Ping service
// @Tags		health
// @ID			ping
// @Produce	json
// @Success	200	{object}	pingResp	"Ping success"
// @Failure	500	{string}	string		"Internal server error"
// @Router		/ping [get].
func Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	resp := pingResp{
		Timestamp: time.Now().String(),
	}

	if err := httputil.WriteJSON(ctx, w, http.StatusOK, resp); err != nil {
		log.Error(ctx, "ping failed", zap.Error(err))

		return
	}

	log.Info(ctx, "ping success")
}

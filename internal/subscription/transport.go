package subscription

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/yushafro/effective-mobile-tz/docs"
	"github.com/yushafro/effective-mobile-tz/pkg/date"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/middleware"
	"go.uber.org/zap"
)

const (
	getPattern     = "GET /subscriptions/{id}"
	createPattern  = "POST /subscriptions"
	updatePattern  = "PUT /subscriptions/{id}"
	deletePattern  = "DELETE /subscriptions/{id}"
	listPattern    = "GET /subscriptions"
	sumPattern     = "GET /subscriptions/sum"
	swaggerPattern = "/swagger/"
)

type server struct {
	service Service
	server  *http.Server
}

func NewServer(service Service, cfg Config) *server {
	return &server{
		service: service,
		server: &http.Server{
			Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
			ReadTimeout:       cfg.ReadTO,
			WriteTimeout:      cfg.WriteTO,
			IdleTimeout:       cfg.IdleTO,
			ReadHeaderTimeout: cfg.ReadHTO,
			Handler:           nil,
		},
	}
}

func (s *server) Start() error {
	return s.server.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *server) RegisterHandlers(logCfg logger.Config) {
	mux := http.NewServeMux()

	mux.HandleFunc(swaggerPattern, httpSwagger.Handler())
	mux.Handle(
		getPattern,
		middleware.NoBody(http.HandlerFunc(s.get)),
	)
	mux.Handle(
		createPattern,
		middleware.RequireBody(http.HandlerFunc(s.create)),
	)
	mux.Handle(
		updatePattern,
		middleware.RequireBody(http.HandlerFunc(s.update)),
	)
	mux.HandleFunc(deletePattern, s.delete)
	mux.Handle(
		listPattern,
		middleware.NoBody(http.HandlerFunc(s.list)),
	)

	mux.Handle(
		sumPattern,
		middleware.NoBody(http.HandlerFunc(s.sum)),
	)

	s.server.Handler = middleware.Logging(mux, logCfg)
}

func (s *server) validateSub(ctx context.Context, sub *SubReq) error {
	log := logger.FromContext(ctx)

	if sub.ServiceName == "" {
		return errEmptyServiceName
	}
	if sub.Price == 0 {
		return errEmptyPrice
	}
	if sub.UserID == "" {
		return errEmptyUserID
	}
	if sub.StartDate == "" {
		return errEmptyStartDate
	}

	err := uuid.Validate(sub.UserID)
	if err != nil {
		log.Error(ctx, errInvalidUserID.Error(), zap.Error(err))

		return fmt.Errorf("%w: %w", errInvalidUserID, err)
	}

	sub.StartDate, err = date.FormatDateToPGDate(ctx, sub.StartDate)
	if err != nil {
		return err
	}

	if sub.EndDate != nil {
		*sub.EndDate, err = date.FormatDateToPGDate(ctx, *sub.EndDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) validateSubSum(ctx context.Context, subSum *SubSumReq) error {
	if subSum.startDate == "" {
		return errEmptyStartDate
	}
	if subSum.endDate == "" {
		return errEmptyEndDate
	}

	var err error
	if subSum.userID != "" {
		err = uuid.Validate(subSum.userID)
		if err != nil {
			return fmt.Errorf("%w: %w", errInvalidUserID, err)
		}
	}

	subSum.startDate, err = date.FormatDateToPGDate(ctx, subSum.startDate)
	if err != nil {
		return err
	}
	subSum.endDate, err = date.FormatDateToPGDate(ctx, subSum.endDate)
	if err != nil {
		return err
	}

	return nil
}

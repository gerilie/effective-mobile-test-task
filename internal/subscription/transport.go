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
	sumPattern     = "GET /subscriptions/sum"
	swaggerPattern = "/swagger/"
)

type server struct {
	service repository
	server  *http.Server
}

func NewServer(service *service, cfg Config) *server {
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

func (s *server) RegisterHandlers() {
	mux := http.NewServeMux()

	mux.HandleFunc(swaggerPattern, httpSwagger.Handler())
	mux.HandleFunc(getPattern, s.get)
	mux.HandleFunc(deletePattern, s.delete)

	mux.Handle(
		createPattern,
		middleware.CancelingEmptyBody(http.HandlerFunc(s.create)),
	)
	mux.Handle(
		updatePattern,
		middleware.CancelingEmptyBody(http.HandlerFunc(s.update)),
	)

	mux.Handle(
		sumPattern,
		middleware.CancelingEmptyBody(http.HandlerFunc(s.sum)),
	)

	s.server.Handler = middleware.Logging(mux)
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

	err = date.FormatM_Y(ctx, &sub.StartDate)
	if err != nil {
		return err
	}
	err = date.FormatM_Y(ctx, sub.EndDate)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) validateSubSum(ctx context.Context, subSum *SubSumReq) error {
	log := logger.FromContext(ctx)

	if subSum.ServiceName == nil && subSum.UserID == nil {
		return errEmptyServiceUser
	}
	if subSum.StartDate == "" {
		return errEmptyStartDate
	}
	if subSum.EndDate == "" {
		return errEmptyEndDate
	}

	if subSum.UserID != nil {
		err := uuid.Validate(*subSum.UserID)
		if err != nil {
			log.Error(ctx, errInvalidUserID.Error(), zap.Error(err))

			return fmt.Errorf("%w: %w", errInvalidUserID, err)
		}
	}

	err := date.FormatM_Y(ctx, &subSum.StartDate)
	if err != nil {
		return err
	}
	err = date.FormatM_Y(ctx, &subSum.EndDate)
	if err != nil {
		return err
	}

	return nil
}

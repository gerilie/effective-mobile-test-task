package subscription

import (
	"context"
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/yushafro/effective-mobile-tz/docs"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/middleware"
	"github.com/yushafro/effective-mobile-tz/pkg/ping"
	"github.com/yushafro/effective-mobile-tz/pkg/ratelimiter"
	"golang.org/x/time/rate"
)

const (
	getPattern     = "GET /subscriptions/{id}"
	createPattern  = "POST /subscriptions"
	updatePattern  = "PATCH /subscriptions/{id}"
	deletePattern  = "DELETE /subscriptions/{id}"
	listPattern    = "GET /subscriptions"
	sumPattern     = "GET /subscriptions/sum"
	swaggerPattern = "/swagger/"
	pingPattern    = "GET /ping"
)

type server struct {
	service  Service
	server   *http.Server
	validate *validator.Validate
	limiter  ratelimiter.IPRateLimiter
	logger   logger.Logger
}

func NewServer(service Service, cfg Config, log logger.Logger) *server {
	validate := validator.New(validator.WithRequiredStructEnabled())
	limiter := ratelimiter.NewIPRateLimiter(
		rate.Limit(cfg.RLRequestsPerSecond),
		cfg.RLBurst,
	)

	s := &server{
		service:  service,
		validate: validate,
		limiter:  limiter,
		logger:   log,
	}

	s.server = &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		ReadTimeout:       cfg.ReadTO,
		WriteTimeout:      cfg.WriteTO,
		IdleTimeout:       cfg.IdleTO,
		ReadHeaderTimeout: cfg.ReadHTO,
		Handler:           s.buildHandler(),
	}

	return s
}

func (s *server) Start() error {
	return s.server.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *server) buildHandler() http.Handler {
	mux := s.buildRouter()

	return s.buildMiddlewareChain(mux)
}

func (s *server) buildRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(pingPattern, ping.Ping)
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

	return mux
}

func (s *server) buildMiddlewareChain(handler http.Handler) http.Handler {
	handler = middleware.RateLimiter(handler, s.limiter)
	handler = middleware.Logging(handler, s.logger)

	return handler
}

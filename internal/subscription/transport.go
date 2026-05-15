package subscription

import (
	"context"
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
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

// Server defines a HTTP server for handling subscription requests.
type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

type server struct {
	service  Service
	server   *http.Server
	validate *validator.Validate
	limiter  ratelimiter.IPRateLimiter
	logger   logger.Logger
}

// NewServer creates a new HTTP server for handling subscription requests.
func NewServer(service Service, cfg Config, log logger.Logger) Server {
	validate := validator.New(validator.WithRequiredStructEnabled())
	limiter := ratelimiter.NewIPRateLimiter(
		ratelimiter.Config{
			R:               rate.Limit(cfg.RLRequestsPerSecond),
			B:               cfg.RLBurst,
			CleanUpInterval: cfg.RLCleanUpInterval,
			CleanUpMaxIdle:  cfg.RLCLeanUpMaxIdle,
		},
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

// Start starts the HTTP server.
func (s *server) Start() error {
	return s.server.ListenAndServe()
}

// Stop stops the HTTP server.
func (s *server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// buildHandler builds the HTTP handler for the server.
func (s *server) buildHandler() http.Handler {
	mux := s.buildRouter()

	return s.buildMiddlewareChain(mux)
}

// buildRouter builds the HTTP router for the server.
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

// buildMiddlewareChain builds the middleware chain for the server.
func (s *server) buildMiddlewareChain(handler http.Handler) http.Handler {
	handler = middleware.RateLimiter(handler, s.limiter)
	handler = middleware.Logging(handler, s.logger)

	return handler
}

package main

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/yushafro/effective-mobile-tz/internal/config"
	"github.com/yushafro/effective-mobile-tz/internal/subscription"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/postgres"
	"go.uber.org/zap"
)

type server struct {
	cfg     *config.Config
	wg      *sync.WaitGroup
	shSrvCh <-chan struct{}
	errCh   chan<- error
}

func initServer(ctx context.Context, s *server) {
	log := logger.FromContext(ctx)

	db, err := postgres.New(ctx, s.cfg.Postgres)
	if err != nil {
		s.errCh <- err

		return
	}
	defer db.Close()

	repo := subscription.NewPGRepository(db)
	service := subscription.NewService(repo)
	server := subscription.NewServer(service, s.cfg.Subscription)
	server.RegisterHandlers(s.cfg.Logger)

	go func() {
		<-s.shSrvCh
		log.Info(ctx, "server stopping")

		err := server.Stop(ctx)
		if err != nil {
			log.Error(ctx, "Error stopping server", zap.Error(err))
			s.errCh <- err

			return
		}

		log.Info(ctx, "server stopped")
		s.wg.Done()
	}()

	log.Info(ctx, "server started")
	err = server.Start()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(ctx, "Error starting server", zap.Error(err))
		s.wg.Done()
		s.errCh <- err

		return
	}
}

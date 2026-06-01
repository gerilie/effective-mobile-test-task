// Package main provides the entry point for the subscription service.
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yushafro/effective-mobile-tz/internal/config"
	"github.com/yushafro/effective-mobile-tz/pkg/deferfunc"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"

	_ "github.com/yushafro/effective-mobile-tz/docs"
)

const (
	servers = 1
	shTO    = 30 * time.Second
)

//	@title			Subscription API
//	@version		1.0
//	@description	Subscription service

//	@BasePath	/

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.Logger, cfg.Subscription.Env)
	if err != nil {
		panic(err)
	}

	ctx, stop := context.WithTimeout(
		logger.WithLogger(context.Background(), log),
		shTO,
	)

	defer deferfunc.Close(ctx, log.Stop, "stop logger")
	defer stop()

	wg := &sync.WaitGroup{}
	shSrvCh := make(chan struct{})
	errCh := make(chan error)
	go func() {
		if err := <-errCh; err != nil {
			panic(err)
		}
	}()

	wg.Add(servers)
	go initServer(ctx, &server{
		cfg:     cfg,
		wg:      wg,
		shSrvCh: shSrvCh,
		errCh:   errCh,
	})

	shCh := make(chan os.Signal, 1)
	signal.Notify(shCh, os.Interrupt, syscall.SIGINT)
	<-shCh
	log.Info("shutdown signal received")

	close(shSrvCh)
	wg.Wait()
	log.Info("servers stopped")
}

package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yushafro/effective-mobile-tz/pkg/deferfunc"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
)

const (
	servers = 1
	shTO    = 30 * time.Second
)

//	@title			Subscription API
//	@version		1.0
//	@description	This is a service for subscriptions

//	@BasePath	/

func main() {
	bootstrapLog := logger.NewBootstrap()
	ctx, stop := context.WithTimeout(logger.WithLogger(context.Background(), bootstrapLog), shTO)
	defer deferfunc.Close(ctx, bootstrapLog.Stop, "Error stopping logger")
	defer stop()

	cfg, err := initConfig(ctx)
	if err != nil {
		panic(err)
	}

	mainLog := logger.NewWithConfig(cfg.Logger)
	ctx = logger.WithLogger(ctx, mainLog)
	defer deferfunc.Close(ctx, mainLog.Stop, "Error stopping logger")

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
	mainLog.Info(ctx, "shutdown signal received")

	close(shSrvCh)
	wg.Wait()
	mainLog.Info(ctx, "all servers stopped")
}

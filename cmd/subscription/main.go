package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/pflag"
	"github.com/yushafro/effective-mobile-tz/pkg/deferfunc"
	"github.com/yushafro/effective-mobile-tz/pkg/env"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
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
	environment := pflag.String("env", env.Dev, "Environment")

	bootstrapLog := logger.NewBootstrap(*environment)
	envCtx := context.WithValue(context.Background(), env.EnvKey, *environment)

	ctx, stop := context.WithTimeout(
		logger.WithLogger(envCtx, bootstrapLog),
		shTO,
	)
	defer deferfunc.Close(ctx, bootstrapLog.Stop, "stop logger")
	defer stop()

	cfg, err := initConfig(ctx)
	if err != nil {
		bootstrapLog.Error(ctx, "init config", zap.Error(err))
		panic(err)
	}

	mainLog := logger.NewWithConfig(cfg.Logger)
	ctx = logger.WithLogger(ctx, mainLog)
	defer deferfunc.Close(ctx, mainLog.Stop, "stop logger")

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
	mainLog.Info(ctx, "servers stopped")
}

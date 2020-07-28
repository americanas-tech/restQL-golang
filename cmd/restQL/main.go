package main

import (
	"context"
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/internal/platform/web"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	if err := start(); err != nil {
		fmt.Printf("[ERROR] failed to start api : %v", err)
		os.Exit(1)
	}
}

var build string

func start() error {
	//// =========================================================================
	//// Config
	startupStart := time.Now()

	cfg, err := conf.Load(build)
	if err != nil {
		return err
	}

	if cfg.Http.Server.EnableFullPprof {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
	}
	log := logger.New(os.Stdout, logger.LogOptions{
		Enable:    cfg.Logging.Enable,
		TimestampFieldName: cfg.Logging.TimestampFieldName,
		TimeFieldFormat: cfg.Logging.TimeFieldFormat,
		Level:     cfg.Logging.Level,
		Format:    cfg.Logging.Format,
	})
	//// =========================================================================
	//// Start API
	log.Info("initializing api")

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	serverCfg := cfg.Http.Server
	apiHandler, err := web.API(log, cfg)
	if err != nil {
		return err
	}

	api := &fasthttp.Server{
		Name:                          "restql",
		Handler:                       apiHandler,
		TCPKeepalive:                  false,
		ReadTimeout:                   serverCfg.ReadTimeout,
		DisableHeaderNamesNormalizing: true,
	}
	health := &fasthttp.Server{
		Name:                          "health",
		Handler:                       web.Health(log, cfg),
		TCPKeepalive:                  false,
		ReadTimeout:                   serverCfg.ReadTimeout,
		DisableHeaderNamesNormalizing: true,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("api listing", "port", serverCfg.ApiAddr)
		serverErrors <- api.ListenAndServe(":" + serverCfg.ApiAddr)
	}()

	go func() {
		defer log.Info("stopping health")
		log.Info("api health listing", "port", serverCfg.ApiHealthAddr)
		serverErrors <- health.ListenAndServe(":" + serverCfg.ApiHealthAddr)
	}()

	if serverCfg.EnablePprof {
		debug := &fasthttp.Server{Name: "debug", Handler: web.Debug(log, cfg)}
		go func() {
			log.Info("api debug listing", "port", serverCfg.PropfAddr)
			serverErrors <- debug.ListenAndServe(":" + serverCfg.PropfAddr)
		}()
	}

	//// =========================================================================
	//// Shutdown
	startupDelay := time.Since(startupStart)
	log.Info("application running", "startup-time", startupDelay.String())

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")
	case sig := <-shutdownSignal:
		log.Info("starting shutdown", "signal", sig)

		timeout, cancel := context.WithTimeout(context.Background(), serverCfg.GracefulShutdownTimeout)
		defer cancel()
		err := shutdown(timeout, log, api, health)

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}

func shutdown(ctx context.Context, log *logger.Logger, servers ...*fasthttp.Server) error {
	var groupErr error
	var g errgroup.Group
	done := make(chan struct{})

	go func() {
		groupErr = g.Wait()
		done <- struct{}{}
	}()

	for _, s := range servers {
		s := s
		g.Go(func() error {
			log.Debug("starting shutdown", "server", s.Name)
			err := s.Shutdown()
			if err != nil {
				log.Error(fmt.Sprintf("%s graceful shutdown did not complete", s.Name), err)
			}
			log.Debug("shutdown finished", "server", s.Name)
			return err
		})
	}

	select {
	case <-ctx.Done():
		return errors.New("graceful shutdown did not complete")
	case <-done:
		return groupErr
	}
}

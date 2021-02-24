package cmd

import (
	"context"
	"fmt"
	"github.com/b2wdigital/restQL-golang/v5/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v5/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/v5/internal/platform/web"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
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

var build string

// Start initialize a restQL runtime as a server
func Start() {
	if err := startServer(); err != nil {
		fmt.Printf("[ERROR] failed to start restQL : %v", err)
		os.Exit(1)
	}
}

func startServer() error {
	//// =========================================================================
	//// Config
	startupStart := time.Now()

	cfg, err := conf.Load(build)
	if err != nil {
		return err
	}

	if cfg.HTTP.Server.EnableFullPprof {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
	}
	log := logger.New(os.Stdout, logger.LogOptions{
		Enable:               cfg.Logging.Enable,
		TimestampFieldName:   cfg.Logging.TimestampFieldName,
		TimestampFieldFormat: cfg.Logging.TimestampFieldFormat,
		Level:                cfg.Logging.Level,
		Format:               cfg.Logging.Format,
	})
	//// =========================================================================
	//// Start API
	log.Info("initializing api")

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	serverCfg := cfg.HTTP.Server
	apiHandler, err := web.API(log, cfg)
	if err != nil {
		return err
	}

	api := &fasthttp.Server{
		Name:                          "restql",
		Handler:                       apiHandler,
		TCPKeepalive:                  true,
		IdleTimeout:                   serverCfg.IdleTimeout,
		ReadTimeout:                   serverCfg.ReadTimeout,
		DisableHeaderNamesNormalizing: true,
	}
	health := &fasthttp.Server{
		Name:                          "health",
		Handler:                       web.Health(log, cfg),
		TCPKeepalive:                  true,
		IdleTimeout:                   serverCfg.IdleTimeout,
		ReadTimeout:                   serverCfg.ReadTimeout,
		DisableHeaderNamesNormalizing: true,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("api listing", "port", serverCfg.APIAddr)
		serverErrors <- api.ListenAndServe(":" + serverCfg.APIAddr)
	}()

	go func() {
		defer log.Info("stopping health")
		log.Info("api health listing", "port", serverCfg.APIHealthAddr)
		serverErrors <- health.ListenAndServe(":" + serverCfg.APIHealthAddr)
	}()

	if serverCfg.EnablePprof {
		debug := &fasthttp.Server{Name: "debug", Handler: web.Debug(log)}
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

func shutdown(ctx context.Context, log restql.Logger, servers ...*fasthttp.Server) error {
	var groupErr error
	var g errgroup.Group
	done := make(chan struct{})

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

	go func() {
		groupErr = g.Wait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.New("graceful shutdown did not complete")
	case <-done:
		return groupErr
	}
}

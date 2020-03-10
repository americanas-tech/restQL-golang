package main

import (
	"context"
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/b2wdigital/restQL-golang/internal/plataform/logger"
	"github.com/b2wdigital/restQL-golang/internal/plataform/web"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
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
	//// Configuration
	config := conf.New(build)
	log := logger.New(os.Stdout, config)

	startupConf, err := newStartupConfig(config, log)
	if err != nil {
		return err
	}

	//// =========================================================================
	//// Start API
	log.Info("initializing api")

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	api := &fasthttp.Server{
		Name:         "api",
		Handler:      web.API(config, log),
		TCPKeepalive: false,
		ReadTimeout:  startupConf.ReadTimeout,
	}
	health := &fasthttp.Server{
		Name:         "health",
		Handler:      web.Health(config, log),
		TCPKeepalive: false,
		ReadTimeout:  startupConf.ReadTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("api listing", "port", startupConf.ApiAddr)
		serverErrors <- api.ListenAndServe(startupConf.ApiAddr)
	}()

	go func() {
		defer log.Info("stopping health")
		log.Info("api health listing", "port", startupConf.ApiHealthAddr)
		serverErrors <- health.ListenAndServe(startupConf.ApiHealthAddr)
	}()

	if startupConf.Env == "development" {
		debug := &fasthttp.Server{Name: "debug", Handler: web.Debug(config, log)}
		go func() {
			log.Info("api debug listing", "port", startupConf.DebugAddr)
			serverErrors <- debug.ListenAndServe(startupConf.DebugAddr)
		}()
	}

	//// =========================================================================
	//// Shutdown
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")
	case sig := <-shutdownSignal:
		log.Info("starting shutdown", "signal", sig)

		timeout, cancel := context.WithTimeout(context.Background(), startupConf.GracefulShutdownTimeout)
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

type startupConfig struct {
	Env                     string
	ApiAddr                 string
	ApiHealthAddr           string
	DebugAddr               string
	GracefulShutdownTimeout time.Duration
	ReadTimeout             time.Duration
}

func newStartupConfig(config conf.Config, log *logger.Logger) (startupConfig, error) {
	env := config.Env().GetString("ENV")

	apiAddr := config.Env().GetString("PORT")
	if apiAddr == "" {
		return startupConfig{}, errors.New("no http port configured, please set PORT environment variable")
	}
	apiAddr = ":" + apiAddr

	apiHealthAddr := config.Env().GetString("HEALTH_PORT")
	if apiHealthAddr == "" {
		return startupConfig{}, errors.New("no http port configured, please set HEALTH_PORT environment variable")
	}
	apiHealthAddr = ":" + apiHealthAddr

	debugAddr := config.Env().GetString("DEBUG_PORT")
	if debugAddr == "" && env == "development" {
		return startupConfig{}, errors.New("no http port configured, please set DEBUG_PORT environment variable")
	}
	debugAddr = ":" + debugAddr

	startupConf := startupConfig{
		Env:                     env,
		ApiAddr:                 apiAddr,
		ApiHealthAddr:           apiHealthAddr,
		DebugAddr:               debugAddr,
		GracefulShutdownTimeout: 5 * time.Second,
		ReadTimeout:             2 * time.Second,
	}

	setWebTimeouts(config, &startupConf, log)

	return startupConf, nil
}

func setWebTimeouts(config conf.Config, startupConf *startupConfig, log *logger.Logger) {
	fileConf := struct {
		Web struct {
			GracefulShutdownTimeout string `yaml:"gracefulShutdownTimeout"`
			ReadTimeout             string `yaml:"readTimeout"`
		} `yaml:"web"`
	}{}
	err := config.File().Unmarshal(&fileConf)
	if err != nil {
		log.Error("failed to load file configuration on startup config parsing", err)
		return
	}

	gracefulShutdownTimeout := fileConf.Web.GracefulShutdownTimeout
	if gracefulShutdownTimeout != "" {
		if gst, err := time.ParseDuration(gracefulShutdownTimeout); err != nil {
			log.Error("graceful shutdown timeout parsing failed", err)
		} else {
			startupConf.GracefulShutdownTimeout = gst
		}
	}

	readTimeout := fileConf.Web.ReadTimeout
	if readTimeout != "" {
		if rt, err := time.ParseDuration(readTimeout); err != nil {
			log.Error("read timeout parsing failed", err)
		} else {
			startupConf.ReadTimeout = rt
		}
	}
}

package main

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/b2wdigital/restQL-golang/internal/plataform/logger"
	"github.com/b2wdigital/restQL-golang/internal/plataform/web"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := start(); err != nil {
		fmt.Printf("[ERROR] failed to start api : %v", err)
		os.Exit(1)
	}
}

func start() error {
	//// =========================================================================
	//// Configuration
	config := conf.New()
	startupConf, err := newStartupConfig(config)
	if err != nil {
		return err
	}

	log := logger.New(os.Stdout, config)

	//// =========================================================================
	//// Start API
	log.Info("initializing api")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := fasthttp.Server{Handler: web.API(config, log)}
	health := fasthttp.Server{Handler: web.Health(config, log)}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("api listing", "port", startupConf.ApiAddr)
		serverErrors <- api.ListenAndServe(startupConf.ApiAddr)
	}()

	go func() {
		log.Info("api health listing", "port", startupConf.ApiHealthAddr)
		serverErrors <- health.ListenAndServe(startupConf.ApiHealthAddr)
	}()

	if startupConf.Env == "development" {
		debug := fasthttp.Server{Handler: web.Debug(config, log)}
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
	case sig := <-shutdown:
		log.Info("starting shutdown", "signal", sig)

		err := api.Shutdown()
		if err != nil {
			log.Error("api graceful shutdown did not complete", err)
		}

		err = health.Shutdown()
		if err != nil {
			log.Error("api health graceful shutdown did not complete", err)
		}

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}

type startupConfig struct {
	Env           string
	ApiAddr       string
	ApiHealthAddr string
	DebugAddr     string
}

func newStartupConfig(config conf.Config) (startupConfig, error) {
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
		Env:           env,
		ApiAddr:       apiAddr,
		ApiHealthAddr: apiHealthAddr,
		DebugAddr:     debugAddr,
	}
	return startupConf, nil
}

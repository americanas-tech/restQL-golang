package main

import (
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/b2wdigital/restQL-golang/internal/plataform/web"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := start(); err != nil {
		log.Printf("[ERROR] failed to start api : %v", err)
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

	//// =========================================================================
	//// Start API
	log.Println("[INFO] initializing api")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := fasthttp.Server{Handler: web.API(config)}
	health := fasthttp.Server{Handler: web.Health(config)}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("[INFO] api listing on %s", startupConf.ApiAddr)
		serverErrors <- api.ListenAndServe(startupConf.ApiAddr)
	}()

	go func() {
		log.Printf("[INFO] api health listing on %s", startupConf.ApiHealthAddr)
		serverErrors <- health.ListenAndServe(startupConf.ApiHealthAddr)
	}()

	if startupConf.Env == "development" {
		debug := fasthttp.Server{Handler: web.Debug(config)}
		go func() {
			log.Printf("[INFO] api debug listing on %s", startupConf.DebugAddr)
			serverErrors <- debug.ListenAndServe(startupConf.DebugAddr)
		}()
	}

	//// =========================================================================
	//// Shutdown
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")
	case sig := <-shutdown:
		log.Printf("[INFO] starting shutdown : %v", sig)

		err := api.Shutdown()
		if err != nil {
			log.Printf("[WARN] api graceful shutdown did not complete : %v", err)
		}

		err = health.Shutdown()
		if err != nil {
			log.Printf("[WARN] api health graceful shutdown did not complete : %v", err)
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

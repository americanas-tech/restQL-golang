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
		log.Printf("[ERROR] failed to start api due to %v", err)
		os.Exit(1)
	}
}

func start() error {
	//// =========================================================================
	//// Configuration
	config := conf.New()
	port := ":" + config.Env().GetString("PORT")

	//// =========================================================================
	//// Start API
	log.Println("[INFO] initializing api")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := fasthttp.Server{Handler: web.New(config)}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("[INFO] api listing on %s", port)
		serverErrors <- api.ListenAndServe(port)
	}()

	//// =========================================================================
	//// Shutdown
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")
	case sig := <-shutdown:
		log.Printf("[INFO] starting shutdown : %v", sig)

		err := api.Shutdown()
		if err != nil {
			log.Printf("[WARN] graceful shutdown did not complete : %v", err)
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

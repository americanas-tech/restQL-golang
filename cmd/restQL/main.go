package main

import (
	"github.com/b2wdigital/restQL-golang/internal/plataform/handlers"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := start(); err != nil {
		log.Printf("[ERROR] failed to start api due to %v", err)
		os.Exit(1)
	}
}

func start() error {
	//// =========================================================================
	//// Start API
	log.Println("[INFO] initializing api")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := fasthttp.Server{
		Handler:      handlers.New(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("[INFO] api listing on %s", ":9000")
		serverErrors <- api.ListenAndServe(":9000")
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
			log.Printf("[WARN] graceful shutdown did not complete in %d : %v", 5000, err)
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

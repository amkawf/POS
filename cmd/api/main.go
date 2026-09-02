package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pos-backend/internal/bootstrap"
	"pos-backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	app, err := bootstrap.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:    ":" + app.Config.AppPort,
		Handler: app.Router,
	}

	go func() {
		log.Printf(
			"API server listening on %s",
			server.Addr,
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	if app.Database != nil {
		app.Database.Close()
	}

	log.Println("server stopped")
}

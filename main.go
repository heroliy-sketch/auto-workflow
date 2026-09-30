package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heroliy-sketch/auto-workflow/config"
	"github.com/heroliy-sketch/auto-workflow/internal/httpapi"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)

	cfg, err := config.Load(os.Getenv("AUTO_WORKFLOW_CONFIG"))
	if err != nil {
		log.Fatalf("load configuration: %v", err)
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.APIPort),
		Handler:           httpapi.NewHandler(cfg.AppName),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("starting %s on port %d", cfg.AppName, cfg.APIPort)
		serverErr <- server.ListenAndServe()
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server stopped: %v", err)
		}
	case <-shutdownCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown: %v", err)
		}
	}
}

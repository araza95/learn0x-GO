package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/araza95/learn0x-GO/internal/config"
	"github.com/araza95/learn0x-GO/internal/http/handlers/student"
	"github.com/araza95/learn0x-GO/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database config
	storage, err := sqlite.New(cfg)
	if err != nil {
		slog.Error("Database not connected!")
		log.Fatal(err)
	}

	slog.Info("Database connected...", slog.String("env", cfg.Env))

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	// setup server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	slog.Info("Server started.", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Unable to start server")
		}
	}()
	<-done

	slog.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server!", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown gracefully")
}

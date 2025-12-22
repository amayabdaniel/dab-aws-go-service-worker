package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/database"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/queue"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/repository"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/scheduler"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/worker"
	"github.com/amayabdaniel/dab-aws-go-service-worker/pkg/config"
	"github.com/amayabdaniel/dab-aws-go-service-worker/pkg/logger"
)

func main() {
	cfg := config.Load()
	slog := logger.New(cfg.LogLevel)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	sqsClient, err := queue.NewSQSClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create SQS client: %v", err)
	}

	// Create repository for scheduler
	repo := repository.NewJobRepository(db)
	
	processor := worker.NewProcessor(db, sqsClient, slog)
	scheduler := scheduler.New(repo, sqsClient, slog)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker
	go func() {
		slog.Info("Starting worker...")
		if err := processor.Start(ctx); err != nil {
			slog.Error("Worker error", "error", err)
		}
	}()

	// Start scheduler
	go func() {
		slog.Info("Starting scheduler...")
		if err := scheduler.Start(ctx); err != nil {
			slog.Error("Scheduler error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down worker and scheduler...")
	cancel()
	slog.Info("Worker exited")
}
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/api/handlers"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/database"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/queue"
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

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	sqsClient, err := queue.NewSQSClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create SQS client: %v", err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	h := handlers.New(db, sqsClient, slog)

	router.GET("/health", h.Health)

	// API routes with /api prefix for ALB routing
	api := router.Group("/api")
	{
		api.GET("/health", h.Health)
		api.POST("/jobs", h.CreateJob)
		api.GET("/jobs/:id", h.GetJob)
		api.GET("/jobs", h.ListJobs)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		slog.Info("Starting API server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	slog.Info("Server exited")
}
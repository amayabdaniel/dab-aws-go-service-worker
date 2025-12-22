package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/interfaces"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/models"
)

type Scheduler struct {
	cron   *cron.Cron
	repo   interfaces.Repository
	queue  interfaces.Queue
	logger *slog.Logger
}

func New(repo interfaces.Repository, queue interfaces.Queue, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		cron:   cron.New(cron.WithSeconds()),
		repo:   repo,
		queue:  queue,
		logger: logger,
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.logger.Info("Starting scheduler")

	// Every 5 minutes: Clean up old completed jobs
	_, err := s.cron.AddFunc("0 */5 * * * *", func() {
		s.cleanupOldJobs(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed to add cleanup job: %w", err)
	}

	// Every hour: Generate system health report
	_, err = s.cron.AddFunc("0 0 * * * *", func() {
		s.generateHealthReport(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed to add health report job: %w", err)
	}

	// Daily at 2 AM: Data aggregation
	_, err = s.cron.AddFunc("0 0 2 * * *", func() {
		s.performDataAggregation(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed to add aggregation job: %w", err)
	}

	// Every 30 seconds: Process batch imports (if any)
	_, err = s.cron.AddFunc("*/30 * * * * *", func() {
		s.processBatchImports(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed to add batch import job: %w", err)
	}

	s.cron.Start()
	
	// Wait for context cancellation
	<-ctx.Done()
	s.logger.Info("Shutting down scheduler")
	
	ctxStop := s.cron.Stop()
	<-ctxStop.Done()
	
	return nil
}

func (s *Scheduler) cleanupOldJobs(ctx context.Context) {
	s.logger.Info("Running cleanup job")
	
	job := &models.Job{
		ID:     uuid.New(),
		Type:   "cleanup",
		Data:   "Remove completed jobs older than 7 days",
		Status: models.JobStatusPending,
	}
	
	if err := s.repo.CreateJob(job); err != nil {
		s.logger.Error("failed to create cleanup job", "error", err)
		return
	}
	
	if err := s.queue.SendMessage(ctx, job.ID.String()); err != nil {
		s.logger.Error("failed to queue cleanup job", "error", err)
	}
}

func (s *Scheduler) generateHealthReport(ctx context.Context) {
	s.logger.Info("Running health report job")
	
	job := &models.Job{
		ID:     uuid.New(),
		Type:   "health-report",
		Data:   fmt.Sprintf("Generate system health report at %s", time.Now().Format(time.RFC3339)),
		Status: models.JobStatusPending,
	}
	
	if err := s.repo.CreateJob(job); err != nil {
		s.logger.Error("failed to create health report job", "error", err)
		return
	}
	
	if err := s.queue.SendMessage(ctx, job.ID.String()); err != nil {
		s.logger.Error("failed to queue health report job", "error", err)
	}
}

func (s *Scheduler) performDataAggregation(ctx context.Context) {
	s.logger.Info("Running data aggregation job")
	
	job := &models.Job{
		ID:     uuid.New(),
		Type:   "data-aggregation",
		Data:   "Aggregate daily metrics and statistics",
		Status: models.JobStatusPending,
	}
	
	if err := s.repo.CreateJob(job); err != nil {
		s.logger.Error("failed to create aggregation job", "error", err)
		return
	}
	
	if err := s.queue.SendMessage(ctx, job.ID.String()); err != nil {
		s.logger.Error("failed to queue aggregation job", "error", err)
	}
}

func (s *Scheduler) processBatchImports(ctx context.Context) {
	// Check if there are any pending batch import requests
	jobs, err := s.repo.ListJobs("pending", 10)
	if err != nil {
		s.logger.Error("failed to list pending jobs", "error", err)
		return
	}
	
	batchCount := 0
	for _, job := range jobs {
		if job.Type == "batch-import" {
			batchCount++
		}
	}
	
	if batchCount > 0 {
		s.logger.Info("Found batch import jobs to process", "count", batchCount)
	}
}
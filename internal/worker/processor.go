package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/models"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/queue"
)

type Processor struct {
	db     *gorm.DB
	queue  *queue.SQSClient
	logger *slog.Logger
}

func NewProcessor(db *gorm.DB, queue *queue.SQSClient, logger *slog.Logger) *Processor {
	return &Processor{
		db:     db,
		queue:  queue,
		logger: logger,
	}
}

func (p *Processor) Start(ctx context.Context) error {
	p.logger.Info("Worker started")

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Worker shutting down")
			return nil
		default:
			messages, err := p.queue.ReceiveMessages(ctx)
			if err != nil {
				p.logger.Error("failed to receive messages", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, msg := range messages {
				if err := p.processMessage(ctx, msg); err != nil {
					p.logger.Error("failed to process message", "error", err)
				}
			}
		}
	}
}

func (p *Processor) processMessage(ctx context.Context, msg types.Message) error {
	if msg.Body == nil {
		return fmt.Errorf("message body is nil")
	}
	
	var jobMsg queue.JobMessage
	if err := json.Unmarshal([]byte(*msg.Body), &jobMsg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	jobID, err := uuid.Parse(jobMsg.JobID)
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	var job models.Job
	if err := p.db.First(&job, "id = ?", jobID).Error; err != nil {
		return fmt.Errorf("failed to find job: %w", err)
	}

	job.Status = models.JobStatusProcessing
	if err := p.db.Save(&job).Error; err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	p.logger.Info("Processing job", "job_id", job.ID)

	result, err := p.processJob(&job)
	if err != nil {
		job.Status = models.JobStatusFailed
		job.Error = err.Error()
	} else {
		job.Status = models.JobStatusCompleted
		job.Result = result
	}

	if err := p.db.Save(&job).Error; err != nil {
		return fmt.Errorf("failed to update job result: %w", err)
	}

	if msg.ReceiptHandle != nil {
		if err := p.queue.DeleteMessage(ctx, *msg.ReceiptHandle); err != nil {
			return fmt.Errorf("failed to delete message: %w", err)
		}
	}

	p.logger.Info("Job processed", "job_id", job.ID, "status", job.Status)
	return nil
}

func (p *Processor) processJob(job *models.Job) (*models.JobResult, error) {
	p.logger.Info("Processing job", "type", job.Type, "id", job.ID)
	
	startTime := time.Now()
	
	switch job.Type {
	case "cleanup":
		return p.processCleanupJob(job)
	case "health-report":
		return p.processHealthReportJob(job)
	case "data-aggregation":
		return p.processDataAggregationJob(job)
	case "batch-import":
		return p.processBatchImportJob(job)
	case "data-processing":
		return p.processDataJob(job)
	default:
		// Generic processing for unknown types
		time.Sleep(1 * time.Second)
		return &models.JobResult{
			ProcessedAt: time.Now(),
			InputCount:  len(job.Data),
			Message:     fmt.Sprintf("Job of type '%s' processed in %v", job.Type, time.Since(startTime)),
		}, nil
	}
}

func (p *Processor) processCleanupJob(job *models.Job) (*models.JobResult, error) {
	// Delete completed jobs older than 7 days
	cutoffDate := time.Now().AddDate(0, 0, -7)
	
	var deletedCount int64
	err := p.db.Model(&models.Job{}).
		Where("status = ? AND updated_at < ?", models.JobStatusCompleted, cutoffDate).
		Count(&deletedCount).
		Delete(&models.Job{}).Error
		
	if err != nil {
		return nil, fmt.Errorf("failed to cleanup old jobs: %w", err)
	}
	
	return &models.JobResult{
		ProcessedAt: time.Now(),
		InputCount:  int(deletedCount),
		Message:     fmt.Sprintf("Cleaned up %d old completed jobs", deletedCount),
	}, nil
}

func (p *Processor) processHealthReportJob(job *models.Job) (*models.JobResult, error) {
	// Generate health metrics
	var metrics struct {
		TotalJobs      int64
		PendingJobs    int64
		ProcessingJobs int64
		CompletedJobs  int64
		FailedJobs     int64
	}
	
	p.db.Model(&models.Job{}).Count(&metrics.TotalJobs)
	p.db.Model(&models.Job{}).Where("status = ?", models.JobStatusPending).Count(&metrics.PendingJobs)
	p.db.Model(&models.Job{}).Where("status = ?", models.JobStatusProcessing).Count(&metrics.ProcessingJobs)
	p.db.Model(&models.Job{}).Where("status = ?", models.JobStatusCompleted).Count(&metrics.CompletedJobs)
	p.db.Model(&models.Job{}).Where("status = ?", models.JobStatusFailed).Count(&metrics.FailedJobs)
	
	// In production, this would send to CloudWatch or S3
	p.logger.Info("Health Report Generated",
		"total", metrics.TotalJobs,
		"pending", metrics.PendingJobs,
		"processing", metrics.ProcessingJobs,
		"completed", metrics.CompletedJobs,
		"failed", metrics.FailedJobs,
	)
	
	return &models.JobResult{
		ProcessedAt: time.Now(),
		InputCount:  int(metrics.TotalJobs),
		Message:     fmt.Sprintf("Health report: Total=%d, Pending=%d, Processing=%d, Completed=%d, Failed=%d",
			metrics.TotalJobs, metrics.PendingJobs, metrics.ProcessingJobs, metrics.CompletedJobs, metrics.FailedJobs),
	}, nil
}

func (p *Processor) processDataAggregationJob(job *models.Job) (*models.JobResult, error) {
	// Aggregate daily statistics
	yesterday := time.Now().AddDate(0, 0, -1)
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	var dailyStats struct {
		JobsCreated   int64
		JobsCompleted int64
		JobsFailed    int64
		AvgProcessing float64
	}
	
	p.db.Model(&models.Job{}).
		Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).
		Count(&dailyStats.JobsCreated)
		
	p.db.Model(&models.Job{}).
		Where("status = ? AND updated_at BETWEEN ? AND ?", models.JobStatusCompleted, startOfDay, endOfDay).
		Count(&dailyStats.JobsCompleted)
		
	p.db.Model(&models.Job{}).
		Where("status = ? AND updated_at BETWEEN ? AND ?", models.JobStatusFailed, startOfDay, endOfDay).
		Count(&dailyStats.JobsFailed)
	
	// In production, this would store in a metrics table or send to data warehouse
	p.logger.Info("Daily aggregation completed",
		"date", yesterday.Format("2006-01-02"),
		"created", dailyStats.JobsCreated,
		"completed", dailyStats.JobsCompleted,
		"failed", dailyStats.JobsFailed,
	)
	
	return &models.JobResult{
		ProcessedAt: time.Now(),
		InputCount:  int(dailyStats.JobsCreated),
		Message:     fmt.Sprintf("Aggregated stats for %s: Created=%d, Completed=%d, Failed=%d",
			yesterday.Format("2006-01-02"), dailyStats.JobsCreated, dailyStats.JobsCompleted, dailyStats.JobsFailed),
	}, nil
}

func (p *Processor) processBatchImportJob(job *models.Job) (*models.JobResult, error) {
	// Simulate processing batch import
	// In production, this would:
	// 1. Parse CSV/JSON from S3
	// 2. Validate each record
	// 3. Bulk insert into database
	// 4. Generate import report
	
	recordCount := len(job.Data) / 10 // Simulate record count
	time.Sleep(100 * time.Millisecond * time.Duration(recordCount))
	
	return &models.JobResult{
		ProcessedAt: time.Now(),
		InputCount:  recordCount,
		Message:     fmt.Sprintf("Batch import completed: %d records processed", recordCount),
	}, nil
}

func (p *Processor) processDataJob(job *models.Job) (*models.JobResult, error) {
	// Simulate data processing
	// In production, this could involve:
	// - ETL operations
	// - API calls to external services
	// - Complex calculations
	// - Report generation
	
	processingTime := time.Duration(len(job.Data)*10) * time.Millisecond
	time.Sleep(processingTime)
	
	return &models.JobResult{
		ProcessedAt: time.Now(),
		InputCount:  len(job.Data),
		Message:     fmt.Sprintf("Data processed successfully in %v", processingTime),
	}, nil
}
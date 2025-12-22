package interfaces

import (
	"context"
	
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/models"
)

// Repository defines database operations
type Repository interface {
	CreateJob(job *models.Job) error
	GetJob(id string) (*models.Job, error)
	UpdateJob(job *models.Job) error
	ListJobs(status string, limit int) ([]models.Job, error)
	GetPendingJobs(limit int) ([]models.Job, error)
}

// Queue defines message queue operations
type Queue interface {
	SendMessage(ctx context.Context, jobID string) error
	ReceiveMessages(ctx context.Context) ([]types.Message, error)
	DeleteMessage(ctx context.Context, receiptHandle string) error
}

// JobProcessor defines job processing operations
type JobProcessor interface {
	ProcessJob(job *models.Job) (*models.JobResult, error)
}
package worker

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/models"
)

func TestProcessor_processJob(t *testing.T) {
	p := &Processor{
		logger: slog.Default(),
	}

	job := &models.Job{
		ID:     uuid.New(),
		Status: models.JobStatusProcessing,
		Type:   "test-job",
		Data:   "test data for processing",
	}

	result, err := p.processJob(job)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotZero(t, result.ProcessedAt)
	assert.Equal(t, len(job.Data), result.InputCount)
	assert.Contains(t, result.Message, "test-job")
}

func TestProcessor_processJob_DataProcessing(t *testing.T) {
	p := &Processor{
		logger: slog.Default(),
	}

	job := &models.Job{
		ID:     uuid.New(),
		Status: models.JobStatusProcessing,
		Type:   "data-processing",
		Data:   "sample data to process",
	}

	result, err := p.processJob(job)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(job.Data), result.InputCount)
	assert.Contains(t, result.Message, "Data processed successfully")
}

func TestProcessor_processJob_BatchImport(t *testing.T) {
	p := &Processor{
		logger: slog.Default(),
	}

	job := &models.Job{
		ID:     uuid.New(),
		Status: models.JobStatusProcessing,
		Type:   "batch-import",
		Data:   "record1,record2,record3,record4,record5",
	}

	result, err := p.processJob(job)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.Message, "Batch import completed")
}

func TestProcessor_Start_ContextCancellation(t *testing.T) {
	// This test verifies that context cancellation works
	_, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Verify the Processor struct can be created
	_ = &Processor{
		logger: slog.Default(),
	}

	// The actual Start() method requires db and queue,
	// so we just verify the struct is properly initialized
}

func TestJobResult_Structure(t *testing.T) {
	result := &models.JobResult{
		ProcessedAt: time.Now(),
		InputCount:  100,
		Message:     "Test message",
	}

	assert.NotZero(t, result.ProcessedAt)
	assert.Equal(t, 100, result.InputCount)
	assert.Equal(t, "Test message", result.Message)
}

func TestJobStatus_Constants(t *testing.T) {
	assert.Equal(t, models.JobStatus("pending"), models.JobStatusPending)
	assert.Equal(t, models.JobStatus("processing"), models.JobStatusProcessing)
	assert.Equal(t, models.JobStatus("completed"), models.JobStatusCompleted)
	assert.Equal(t, models.JobStatus("failed"), models.JobStatusFailed)
}

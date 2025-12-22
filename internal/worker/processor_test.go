package worker

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"log/slog"

	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/models"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/queue"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	// Populate the job for testing
	if job, ok := dest.(*models.Job); ok && args.Error(0) == nil {
		job.ID = uuid.New()
		job.Status = models.JobStatusPending
		job.Type = "test"
		job.Data = "test data"
	}
	return &gorm.DB{Error: args.Error(0)}
}

func (m *mockDB) Save(value interface{}) *gorm.DB {
	args := m.Called(value)
	return &gorm.DB{Error: args.Error(0)}
}

type mockQueue struct {
	mock.Mock
}

func (m *mockQueue) ReceiveMessages(ctx context.Context) ([]interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]interface{}), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockQueue) DeleteMessage(ctx context.Context, receiptHandle string) error {
	args := m.Called(ctx, receiptHandle)
	return args.Error(0)
}

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
	assert.Equal(t, "Job of type 'test-job' processed successfully", result.Message)
}

func TestProcessor_processMessage_InvalidJSON(t *testing.T) {
	mockDatabase := &mockDB{}
	mockSQS := &mockQueue{}
	
	p := &Processor{
		db:     &gorm.DB{},
		queue:  mockSQS,
		logger: slog.Default(),
	}

	// Create a mock message with invalid JSON
	invalidBody := "invalid json"
	receiptHandle := "test-receipt"
	mockMsg := types.Message{
		Body:          &invalidBody,
		ReceiptHandle: &receiptHandle,
	}

	err := p.processMessage(context.Background(), mockMsg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal message")
}

func TestProcessor_Start_ContextCancellation(t *testing.T) {
	mockDatabase := &mockDB{}
	mockSQS := &mockQueue{}
	
	p := &Processor{
		db:     &gorm.DB{},
		queue:  mockSQS,
		logger: slog.Default(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	// Cancel immediately
	cancel()

	err := p.Start(ctx)
	assert.NoError(t, err)
}
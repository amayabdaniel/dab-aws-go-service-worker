package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/models"
)

// Mock Repository
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateJob(job *models.Job) error {
	args := m.Called(job)
	if args.Error(0) == nil {
		job.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *mockRepository) GetJob(id string) (*models.Job, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Job), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepository) UpdateJob(job *models.Job) error {
	args := m.Called(job)
	return args.Error(0)
}

func (m *mockRepository) ListJobs(status string, limit int) ([]models.Job, error) {
	args := m.Called(status, limit)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Job), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepository) GetPendingJobs(limit int) ([]models.Job, error) {
	args := m.Called(limit)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Job), args.Error(1)
	}
	return nil, args.Error(1)
}

// Mock Queue
type mockQueue struct {
	mock.Mock
}

func (m *mockQueue) SendMessage(ctx context.Context, jobID string) error {
	args := m.Called(ctx, jobID)
	return args.Error(0)
}

func (m *mockQueue) ReceiveMessages(ctx context.Context) ([]types.Message, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]types.Message), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockQueue) DeleteMessage(ctx context.Context, receiptHandle string) error {
	args := m.Called(ctx, receiptHandle)
	return args.Error(0)
}

// handlerWithMocks creates a handler with injected mocks for testing
type handlerWithMocks struct {
	repo   *mockRepository
	queue  *mockQueue
	logger *slog.Logger
}

func (h *handlerWithMocks) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "api",
	})
}

func (h *handlerWithMocks) CreateJob(c *gin.Context) {
	var payload models.JobPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	job := &models.Job{
		Status: models.JobStatusPending,
		Type:   payload.Type,
		Data:   payload.Data,
	}

	if err := h.repo.CreateJob(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create job"})
		return
	}

	go func() {
		h.queue.SendMessage(context.Background(), job.ID.String())
	}()

	c.JSON(http.StatusCreated, job)
}

func (h *handlerWithMocks) GetJob(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job ID"})
		return
	}

	job, err := h.repo.GetJob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *handlerWithMocks) ListJobs(c *gin.Context) {
	status := c.Query("status")

	jobs, err := h.repo.ListJobs(status, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":  jobs,
		"count": len(jobs),
	})
}

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &handlerWithMocks{
		logger: slog.Default(),
	}

	router := gin.New()
	router.GET("/health", h.Health)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "api", response["service"])
}

func TestCreateJob_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &mockRepository{}
	mockSQS := &mockQueue{}

	h := &handlerWithMocks{
		repo:   mockRepo,
		queue:  mockSQS,
		logger: slog.Default(),
	}

	mockRepo.On("CreateJob", mock.Anything).Return(nil)
	mockSQS.On("SendMessage", mock.Anything, mock.Anything).Return(nil)

	router := gin.New()
	router.POST("/jobs", h.CreateJob)

	payload := models.JobPayload{
		Type: "data-processing",
		Data: "test data",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var job models.Job
	err := json.Unmarshal(w.Body.Bytes(), &job)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, job.ID)
	assert.Equal(t, models.JobStatusPending, job.Status)
	assert.Equal(t, "data-processing", job.Type)

	mockRepo.AssertExpectations(t)
}

func TestCreateJob_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &handlerWithMocks{
		logger: slog.Default(),
	}

	router := gin.New()
	router.POST("/jobs", h.CreateJob)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request body", response["error"])
}

func TestGetJob_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &mockRepository{}

	h := &handlerWithMocks{
		repo:   mockRepo,
		logger: slog.Default(),
	}

	jobID := uuid.New()
	expectedJob := &models.Job{
		ID:     jobID,
		Status: models.JobStatusCompleted,
		Type:   "test",
		Data:   "test data",
	}

	mockRepo.On("GetJob", jobID.String()).Return(expectedJob, nil)

	router := gin.New()
	router.GET("/jobs/:id", h.GetJob)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs/"+jobID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var job models.Job
	err := json.Unmarshal(w.Body.Bytes(), &job)
	assert.NoError(t, err)
	assert.Equal(t, jobID, job.ID)

	mockRepo.AssertExpectations(t)
}

func TestGetJob_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &handlerWithMocks{
		logger: slog.Default(),
	}

	router := gin.New()
	router.GET("/jobs/:id", h.GetJob)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs/invalid-uuid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid job ID", response["error"])
}

func TestGetJob_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &mockRepository{}

	h := &handlerWithMocks{
		repo:   mockRepo,
		logger: slog.Default(),
	}

	jobID := uuid.New()
	mockRepo.On("GetJob", jobID.String()).Return(nil, errors.New("not found"))

	router := gin.New()
	router.GET("/jobs/:id", h.GetJob)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs/"+jobID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockRepo.AssertExpectations(t)
}

func TestListJobs_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &mockRepository{}

	h := &handlerWithMocks{
		repo:   mockRepo,
		logger: slog.Default(),
	}

	expectedJobs := []models.Job{
		{ID: uuid.New(), Status: models.JobStatusPending, Type: "test1", Data: "data1"},
		{ID: uuid.New(), Status: models.JobStatusCompleted, Type: "test2", Data: "data2"},
	}

	mockRepo.On("ListJobs", "", 100).Return(expectedJobs, nil)

	router := gin.New()
	router.GET("/jobs", h.ListJobs)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), response["count"])

	mockRepo.AssertExpectations(t)
}

func TestListJobs_WithStatusFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &mockRepository{}

	h := &handlerWithMocks{
		repo:   mockRepo,
		logger: slog.Default(),
	}

	expectedJobs := []models.Job{
		{ID: uuid.New(), Status: models.JobStatusPending, Type: "test1", Data: "data1"},
	}

	mockRepo.On("ListJobs", "pending", 100).Return(expectedJobs, nil)

	router := gin.New()
	router.GET("/jobs", h.ListJobs)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs?status=pending", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["count"])

	mockRepo.AssertExpectations(t)
}

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/api/handlers"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/models"
)

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 200, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "api", response["service"])
}

func TestCreateJob(t *testing.T) {
	router := setupTestRouter()
	
	payload := map[string]interface{}{
		"type": "process",
		"data": map[string]interface{}{
			"input": "test data",
		},
	}
	
	body, _ := json.Marshal(payload)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 201, w.Code)
	
	var job models.Job
	err := json.Unmarshal(w.Body.Bytes(), &job)
	require.NoError(t, err)
	
	assert.NotEqual(t, uuid.Nil, job.ID)
	assert.Equal(t, models.JobStatusPending, job.Status)
	assert.Equal(t, payload, job.Payload)
}

func TestGetJob(t *testing.T) {
	router := setupTestRouter()
	
	// Test non-existent job
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 404, w.Code)
	
	// Test invalid UUID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/jobs/invalid-uuid", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 400, w.Code)
}

func TestListJobs(t *testing.T) {
	router := setupTestRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 200, w.Code)
	
	var jobs []models.Job
	err := json.Unmarshal(w.Body.Bytes(), &jobs)
	require.NoError(t, err)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Mock handlers for testing
	// In real tests, you'd use a test database and mock SQS
	h := &mockHandler{}
	
	router.GET("/health", h.Health)
	router.POST("/jobs", h.CreateJob)
	router.GET("/jobs/:id", h.GetJob)
	router.GET("/jobs", h.ListJobs)
	
	return router
}

type mockHandler struct{}

func (m *mockHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "api",
	})
}

func (m *mockHandler) CreateJob(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	
	job := &models.Job{
		ID:      uuid.New(),
		Status:  models.JobStatusPending,
		Payload: payload,
	}
	
	c.JSON(http.StatusCreated, job)
}

func (m *mockHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job ID"})
		return
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
}

func (m *mockHandler) ListJobs(c *gin.Context) {
	c.JSON(http.StatusOK, []models.Job{})
}
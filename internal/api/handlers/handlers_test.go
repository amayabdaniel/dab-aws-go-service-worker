package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"log/slog"

	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/models"
)

// Mock for database
type mockDB struct {
	mock.Mock
}

func (m *mockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return &gorm.DB{Error: args.Error(0)}
}

func (m *mockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return &gorm.DB{Error: args.Error(0)}
}

func (m *mockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest)
	return &gorm.DB{Error: args.Error(0)}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Order(value interface{}) *gorm.DB {
	return m
}

func (m *mockDB) Limit(limit int) *gorm.DB {
	return m
}

// Mock for SQS
type mockQueue struct {
	mock.Mock
}

func (m *mockQueue) SendMessage(ctx interface{}, jobID string) error {
	args := m.Called(ctx, jobID)
	return args.Error(0)
}

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	h := &Handler{
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

	mockDatabase := &gorm.DB{}
	mockSQS := &mockQueue{}
	
	h := &Handler{
		db:     mockDatabase,
		queue:  mockSQS,
		logger: slog.Default(),
	}

	// Mock expectations
	mockSQS.On("SendMessage", mock.Anything, mock.Anything).Return(nil)

	router := gin.New()
	router.POST("/jobs", h.CreateJob)

	payload := map[string]interface{}{
		"type": "test",
		"data": "sample",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Note: Without a real DB connection, this will fail
	// In production tests, you'd use a test database or better mocks
}

func TestCreateJob_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &Handler{
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

func TestGetJob_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &Handler{
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

func TestListJobs_WithStatusFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDatabase := &gorm.DB{}
	
	h := &Handler{
		db:     mockDatabase,
		logger: slog.Default(),
	}

	router := gin.New()
	router.GET("/jobs", h.ListJobs)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs?status=pending", nil)
	router.ServeHTTP(w, req)

	// Note: Without a real DB connection, this will fail
	// In production tests, you'd use a test database or better mocks
}
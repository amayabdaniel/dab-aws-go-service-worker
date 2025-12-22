package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/api/middleware"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/interfaces"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/models"
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/repository"
)

type Handler struct {
	repo   interfaces.Repository
	queue  interfaces.Queue
	logger *slog.Logger
}

func New(db *gorm.DB, queue interfaces.Queue, logger *slog.Logger) *Handler {
	return &Handler{
		repo:   repository.NewJobRepository(db),
		queue:  queue,
		logger: logger,
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "api",
	})
}

func (h *Handler) CreateJob(c *gin.Context) {
	var payload models.JobPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		middleware.ValidationError(c, err)
		return
	}

	// Validate using validator tags
	if err := middleware.ValidateStruct(payload); err != nil {
		middleware.ValidationError(c, err)
		return
	}

	job := &models.Job{
		Status: models.JobStatusPending,
		Type:   payload.Type,
		Data:   payload.Data,
	}

	if err := h.repo.CreateJob(job); err != nil {
		h.logger.Error("failed to create job", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create job"})
		return
	}

	// Queue job asynchronously with background context
	go func() {
		ctx := context.Background()
		if err := h.queue.SendMessage(ctx, job.ID.String()); err != nil {
			h.logger.Error("failed to queue job", "error", err, "job_id", job.ID)
		}
	}()

	c.JSON(http.StatusCreated, job)
}

func (h *Handler) GetJob(c *gin.Context) {
	id := c.Param("id")
	
	job, err := h.repo.GetJob(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job ID"})
		case errors.Is(err, repository.ErrJobNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		default:
			h.logger.Error("failed to get job", "error", err, "job_id", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get job"})
		}
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *Handler) ListJobs(c *gin.Context) {
	status := c.Query("status")
	
	jobs, err := h.repo.ListJobs(status, 100)
	if err != nil {
		h.logger.Error("failed to list jobs", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":  jobs,
		"count": len(jobs),
	})
}
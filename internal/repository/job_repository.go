package repository

import (
	"errors"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
	
	"github.com/amayabdaniel/dab-aws-go-service-worker/internal/models"
)

var (
	ErrJobNotFound = errors.New("job not found")
	ErrInvalidID   = errors.New("invalid job ID")
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) CreateJob(job *models.Job) error {
	if job == nil {
		return errors.New("job cannot be nil")
	}
	return r.db.Create(job).Error
}

func (r *JobRepository) GetJob(id string) (*models.Job, error) {
	jobID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	
	var job models.Job
	if err := r.db.First(&job, "id = ?", jobID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}
	
	return &job, nil
}

func (r *JobRepository) UpdateJob(job *models.Job) error {
	if job == nil {
		return errors.New("job cannot be nil")
	}
	return r.db.Save(job).Error
}

func (r *JobRepository) ListJobs(status string, limit int) ([]models.Job, error) {
	var jobs []models.Job
	
	query := r.db.Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(100) // default limit
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}
	
	return jobs, nil
}

func (r *JobRepository) GetPendingJobs(limit int) ([]models.Job, error) {
	return r.ListJobs(string(models.JobStatusPending), limit)
}
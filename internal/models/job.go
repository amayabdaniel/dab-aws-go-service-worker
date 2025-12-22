package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

type JobPayload struct {
	Type string `json:"type" validate:"required,min=1,max=100"`
	Data string `json:"data" validate:"required,min=1,max=10000"`
}

type JobResult struct {
	ProcessedAt time.Time `json:"processed_at"`
	InputCount  int       `json:"input_count"`
	Message     string    `json:"message"`
}

type Job struct {
	ID        uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primary_key" json:"id"`
	Status    JobStatus   `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Type      string      `gorm:"type:varchar(100);not null" json:"type"`
	Data      string      `gorm:"type:text;not null" json:"data"`
	Result    *JobResult  `gorm:"serializer:json" json:"result,omitempty"`
	Error     string      `gorm:"type:text" json:"error,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func (Job) TableName() string {
	return "jobs"
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}
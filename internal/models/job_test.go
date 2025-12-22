package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJob_BeforeCreate(t *testing.T) {
	job := &Job{
		Status:    JobStatusPending,
		Type:      "test",
		Data:      "test data",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if job.ID != uuid.Nil {
		t.Error("expected nil UUID before BeforeCreate")
	}

	err := job.BeforeCreate(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if job.ID == uuid.Nil {
		t.Error("expected non-nil UUID after BeforeCreate")
	}
}

func TestJobStatus_Values(t *testing.T) {
	statuses := []JobStatus{
		JobStatusPending,
		JobStatusProcessing,
		JobStatusCompleted,
		JobStatusFailed,
	}

	expected := []string{"pending", "processing", "completed", "failed"}

	for i, status := range statuses {
		if string(status) != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], string(status))
		}
	}
}
package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tests := []struct {
		name        string
		databaseURL string
		wantErr     bool
	}{
		{
			name:        "invalid database URL",
			databaseURL: "invalid://url",
			wantErr:     true,
		},
		{
			name:        "empty database URL",
			databaseURL: "",
			wantErr:     true,
		},
		// Note: We can't test actual connection without a real database
		// In a real test environment, you'd use a test database or mock
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Connect(tt.databaseURL)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
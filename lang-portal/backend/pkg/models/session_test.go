package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessionValidation(t *testing.T) {
	now := time.Now()
	later := now.Add(1 * time.Hour)

	testCases := []struct {
		name      string
		session   Session
		expectErr bool
	}{
		{
			name: "Valid Session",
			session: Session{
				ID:         1,
				ActivityID: 1,
				GroupID:    nil,
				StartTime:  now,
				EndTime:    &later,
				Score:      10,
			},
			expectErr: false,
		},
		{
			name: "Invalid Activity ID",
			session: Session{
				ID:         1,
				ActivityID: 0,
				StartTime:  now,
			},
			expectErr: true,
		},
		{
			name: "Invalid Group ID",
			session: Session{
				ID:         1,
				ActivityID: 1,
				GroupID:    intPtr(-1),
				StartTime:  now,
			},
			expectErr: true,
		},
		{
			name: "Zero Start Time",
			session: Session{
				ID:         1,
				ActivityID: 1,
				StartTime:  time.Time{},
			},
			expectErr: true,
		},
		{
			name: "End Time Before Start Time",
			session: Session{
				ID:         1,
				ActivityID: 1,
				StartTime:  now,
				EndTime:    &time.Time{},
			},
			expectErr: true,
		},
		{
			name: "Negative Score",
			session: Session{
				ID:         1,
				ActivityID: 1,
				StartTime:  now,
				Score:      -1,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.session.Validate()
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionDuration(t *testing.T) {
	now := time.Now()
	later := now.Add(30 * time.Minute)

	t.Run("Completed Session", func(t *testing.T) {
		session := Session{
			StartTime: now,
			EndTime:   &later,
		}
		assert.Equal(t, 30*time.Minute, session.Duration())
	})

	t.Run("Incomplete Session", func(t *testing.T) {
		session := Session{
			StartTime: now,
			EndTime:   nil,
		}
		assert.Equal(t, time.Duration(0), session.Duration())
	})
}

func TestSessionIsCompleted(t *testing.T) {
	now := time.Now()
	later := now.Add(30 * time.Minute)

	t.Run("Completed Session", func(t *testing.T) {
		session := Session{
			StartTime: now,
			EndTime:   &later,
		}
		assert.True(t, session.IsCompleted())
	})

	t.Run("Incomplete Session", func(t *testing.T) {
		session := Session{
			StartTime: now,
			EndTime:   nil,
		}
		assert.False(t, session.IsCompleted())
	})
}

// Helper function to create a pointer to an int
func intPtr(i int64) *int64 {
	return &i
}

package models_test

import (
	"testing"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
)

func TestGroup_Validate(t *testing.T) {
	tests := []struct {
		name    string
		group   models.Group
		wantErr bool
	}{
		{
			name: "valid group",
			group: models.Group{
				Name: "Travel Words",
			},
			wantErr: false,
		},
		{
			name: "empty group name",
			group: models.Group{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "too short group name",
			group: models.Group{
				Name: "A",
			},
			wantErr: true,
		},
		{
			name: "too long group name",
			group: models.Group{
				Name: "This is an extremely long group name that exceeds the maximum allowed length",
			},
			wantErr: true,
		},
		{
			name: "group name with invalid characters",
			group: models.Group{
				Name: "Travel Words 123!",
			},
			wantErr: true,
		},
		{
			name: "group name with allowed special characters",
			group: models.Group{
				Name: "Travel-Words_Group",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.group.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Group.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroup_Sanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trim whitespace",
			input:    "  Travel Words  ",
			expected: "Travel Words",
		},
		{
			name:     "no changes needed",
			input:    "Travel Words",
			expected: "Travel Words",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &models.Group{Name: tt.input}
			group.Sanitize()

			if group.Name != tt.expected {
				t.Errorf("Group.Sanitize() = %v, want %v", group.Name, tt.expected)
			}
		})
	}
}

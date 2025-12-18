package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "Age 34 - birthday passes",
			dob:      time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			expected: 34,
		},
		{
			name:     "Age 33 - birthday not yet",
			dob:      time.Date(1991, 12, 25, 0, 0, 0, 0, time.UTC),
			expected: 33,
		},
		{
			name:     "Age 0 - born this year",
			dob:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Age 25",
			dob:      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			age := CalculateAge(tt.dob)

			if age < tt.expected-1 || age > tt.expected+1 {
				t.Errorf("CalculateAge(%v) = %d, want approximately %d", tt.dob, age, tt.expected)
			}
		})
	}
}

func TestCalculateAgeEdgeCases(t *testing.T) {
	now := time.Now()

	dob := now.AddDate(-30, 0, 0)
	age := CalculateAge(dob)
	if age != 30 {
		t.Errorf("Birthday today: expected 30, got %d", age)
	}

	dob = now.AddDate(-30, 0, 1)
	age = CalculateAge(dob)
	if age != 29 {
		t.Errorf("Birthday tomorrow: expected 29, got %d", age)
	}

	dob = now.AddDate(-30, 0, -1)
	age = CalculateAge(dob)
	if age != 30 {
		t.Errorf("Birthday yesterday: expected 30, got %d", age)
	}
}

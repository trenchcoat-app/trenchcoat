package auth_test

import (
	"testing"
	"time"

	"trenchcoat/internal/services/auth"

	"github.com/go-openapi/testify/v2/assert"
)

func TestGetNewSessionExpireTime(t *testing.T) {
	svc := auth.NewAuthService(nil)

	tests := []struct {
		name       string
		offsetSecs int
		tolerance  time.Duration
	}{
		{
			name:       "default session expiry",
			offsetSecs: 86400,
			tolerance:  time.Second,
		},
		{
			name:       "short expiry",
			offsetSecs: 3600,
			tolerance:  time.Second,
		},
		{
			name:       "zero offset",
			offsetSecs: 0,
			tolerance:  time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			before := time.Now()
			result := svc.GetNewSessionExpireTime(tt.offsetSecs)
			after := time.Now()

			expectedMin := before.Add(time.Duration(tt.offsetSecs) * time.Second)
			expectedMax := after.Add(time.Duration(tt.offsetSecs) * time.Second)

			assert.True(t, result.After(expectedMin) || result.Equal(expectedMin),
				"Expected %v to be after or equal to %v", result, expectedMin)
			assert.True(t, result.Before(expectedMax) || result.Equal(expectedMax),
				"Expected %v to be before or equal to %v", result, expectedMax)
		})
	}
}

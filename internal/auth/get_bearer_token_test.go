package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		header        http.Header
		expectError   bool
		expectedToken string
	}{
		{
			name: "valid header",
			header: http.Header{
				"Authorization": []string{"Bearer abc123!"},
			},
			expectError:   false,
			expectedToken: "abc123!",
		},
		{
			name:          "missing authorization header",
			header:        http.Header{},
			expectError:   true,
			expectedToken: "",
		},
		{
			name: "empty authorization header",
			header: http.Header{
				"Authorization": []string{""},
			},
			expectError:   true,
			expectedToken: "",
		},
		{
			name: "trailing whitespace",
			header: http.Header{
				"Authorization": []string{" Bearer abc123!"},
			},
			expectError:   false,
			expectedToken: "abc123!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.header)
			if tt.expectError {
				assert.Error(t, err, "expected no token and error, got '%v', %v", token, err)
			} else {
				assert.Equal(t, tt.expectedToken, token, "expected token %v and no error got '%v', %v", tt.expectedToken, token, err)
			}
		})
	}
}

package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAPIKey(t *testing.T) {

	tests := []struct {
		name        string
		headers     http.Header
		expectError bool
		expect      string
	}{
		{
			name: "valid header",
			headers: http.Header{
				"Authorization": {"ApiKey thisisakey"},
			},
			expectError: false,
			expect:      "thisisakey",
		},
		{
			name: "empty header",
			headers: http.Header{
				"Authorization": {""},
			},
			expectError: true,
			expect:      "",
		},
		{
			name: "malformed header",
			headers: http.Header{
				"Authorization": {"ApiKeythisiskey"},
			},
			expectError: true,
			expect:      "",
		},
		{
			name: "no authorization header",
			headers: http.Header{
				"Content-Type": {"application/json"},
			},
			expectError: true,
			expect:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GetAPIKey(tt.headers)

			if tt.expectError {
				assert.Error(t, err, "expected error and output '%v' == '%v'", key, tt.expect)
			} else {
				assert.NoError(t, err, "expected no error got %v", err)
				assert.Equal(t, tt.expect, key, "key should match")
			}
		})
	}
}

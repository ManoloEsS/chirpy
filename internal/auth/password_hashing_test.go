package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		expectedError bool
	}{
		{
			name:          "valid password",
			password:      "Manolo",
			expectedError: false,
		},
		{
			name:          "empty password",
			password:      "",
			expectedError: false,
		},
		{
			name:          "special chars",
			password:      "P@sswoOrd!#$",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.password, hash)

			}
		})
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "Manolo"
	hash1, err1 := HashPassword(password)
	require.NoError(t, err1)

	hash2, err2 := HashPassword(password)
	require.NoError(t, err2)

	assert.NotEqual(t, hash1, hash2, "same password should produce different hashes")
}

func TestCheckPasswordHash(t *testing.T) {
	password := "Manolo"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name        string
		password    string
		hash        string
		expectMatch bool
	}{
		{
			name:        "correct password",
			password:    password,
			hash:        hash,
			expectMatch: true,
		},
		{
			name:        "wrong password",
			password:    "wrong",
			hash:        hash,
			expectMatch: false,
		},
		{
			name:        "empty password",
			password:    "",
			hash:        hash,
			expectMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, _ := CheckPasswordHash(tt.password, tt.hash)
			assert.Equal(t, tt.expectMatch, matches)
		})
	}

}

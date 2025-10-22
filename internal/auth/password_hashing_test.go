package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)
	tests := []struct {
		name          string
		password      string
		hash          string
		expectedError bool
		matchPassword bool
	}{
		{
			name:          "valid password",
			password:      password1,
			hash:          hash1,
			expectedError: false,
			matchPassword: true,
		},
		{
			name:          "incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			expectedError: false,
			matchPassword: false,
		},
		{
			name:          "password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			expectedError: false,
			matchPassword: false,
		},
		{
			name:          "empty password",
			password:      "",
			hash:          hash1,
			expectedError: false,
			matchPassword: false,
		},
		{
			name:          "invalid hash",
			password:      password1,
			hash:          "invalidhash",
			expectedError: true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.matchPassword, match)
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

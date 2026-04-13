package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key"

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()

	t.Run("valid token is created", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, testSecret, time.Hour)
		require.NoError(t, err)
		assert.NotEmpty(t, tokenString)
	})

	t.Run("token has three parts (header.payload.signature)", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, testSecret, time.Hour)
		require.NoError(t, err)

		parts := strings.Split(tokenString, ".")
		assert.Len(t, parts, 3)
	})
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()

	t.Run("valid token returns correct userID", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, testSecret, time.Hour)
		require.NoError(t, err)

		gotID, err := ValidateJWT(tokenString, testSecret)
		require.NoError(t, err)
		assert.Equal(t, userID, gotID)
	})

	t.Run("wrong secret returns error", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, testSecret, time.Hour)
		require.NoError(t, err)

		_, err = ValidateJWT(tokenString, "wrong-secret")
		assert.Error(t, err)
	})

	t.Run("expired token returns error", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, testSecret, -time.Hour) // already expired
		require.NoError(t, err)

		_, err = ValidateJWT(tokenString, testSecret)
		assert.Error(t, err)
	})

	t.Run("malformed token returns error", func(t *testing.T) {
		_, err := ValidateJWT("not.a.token", testSecret)
		assert.Error(t, err)
	})

	t.Run("empty token returns error", func(t *testing.T) {
		_, err := ValidateJWT("", testSecret)
		assert.Error(t, err)
	})

	t.Run("tampered payload returns error", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, testSecret, time.Hour)
		require.NoError(t, err)

		// Swap out the payload section with a different base64 string
		parts := strings.Split(tokenString, ".")
		parts[1] = "dGFtcGVyZWRwYXlsb2Fk" // base64 of "tamperedpayload"
		tampered := strings.Join(parts, ".")

		_, err = ValidateJWT(tampered, testSecret)
		assert.Error(t, err)
	})

	t.Run("token signed with different method is rejected", func(t *testing.T) {
		// Manually craft a token with alg: none
		unsignedToken := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0" + // {"alg":"none","typ":"JWT"}
			".eyJzdWIiOiIxMjM0NTY3ODkwIn0" + // {"sub":"1234567890"}
			"."

		_, err := ValidateJWT(unsignedToken, testSecret)
		assert.Error(t, err)
	})
}

package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	registeredClaims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)

	secretKey := []byte(tokenSecret)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	myClaims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &myClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("Invalid token")
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to get subject: %w", err)
	}

	userId, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Unable to parse userID: %w", err)
	}
	return userId, err
}

// Package auth provides user authentication with JWTs.
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService manages authentication tokens with generate and validate methods.
type JWTService struct {
	secret string
}

// NewJWTService returns a JWTService
func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: secret}
}

// Claims represents the custom and registered claims of the token.
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken returns a token generated for the given userID.
func (s *JWTService) GenerateToken(userID int) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "rifflog-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a token and returns the Claims associate with it.
func (s *JWTService) ValidateToken(tokenString string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return Claims{}, fmt.Errorf("parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return *claims, nil
	}

	return Claims{}, errors.New("invalid token claims")
}

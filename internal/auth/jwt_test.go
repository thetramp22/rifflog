package auth

import (
	"testing"
)

func TestGenerateAndValidateToken(t *testing.T) {
	service := NewJWTService("test-secret")

	userID := 22

	token, err := service.GenerateToken(userID)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Error validating token: %v", err)
	}

	if claims.UserID != userID {
		t.Fatalf("Expected %v, got %v", userID, claims.UserID)
	}
}

func TestValidateTokenWrongToken(t *testing.T) {
	service1 := NewJWTService("secret-one")
	service2 := NewJWTService("secret-two")

	userID := 22

	token, err := service1.GenerateToken(userID)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}

	_, err = service2.ValidateToken(token)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

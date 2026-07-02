package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/thetramp22/rifflog/internal/models"
)

func TestLogin(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	data := models.LoginRequest{
		Email:    user.User.Email,
		Password: password,
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	req := httptest.NewRequest("POST", "http://localhost:8080/login", bytes.NewReader(jsonBytes))
	w := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("expected 200, got %v, body=%s", status, w.Body.String())
	}

	var got models.LoginResponse
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if got.Token == "" {
		t.Fatal("expected token")
	}

	claims, err := app.JWTService.ValidateToken(got.Token)
	if err != nil {
		t.Fatalf("error validating token: %v", err)
	}

	if diff := cmp.Diff(user.User.ID, claims.UserID); diff != "" {
		t.Errorf("values mismatch (-want +got):\n%s", diff)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	wrongPassword := "wrong"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	data := models.LoginRequest{
		Email:    user.User.Email,
		Password: wrongPassword,
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	req := httptest.NewRequest("POST", "http://localhost:8080/login", bytes.NewReader(jsonBytes))
	w := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusUnauthorized {
		t.Fatalf("expected 200, got %v, body=%s", status, w.Body.String())
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	unknownEmail := "wrongEmail@test.com"
	app, _ := SetupTestUser(t, password)
	defer app.DB.Close()

	data := models.LoginRequest{
		Email:    unknownEmail,
		Password: password,
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	req := httptest.NewRequest("POST", "http://localhost:8080/login", bytes.NewReader(jsonBytes))
	w := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusUnauthorized {
		t.Fatalf("expected 200, got %v, body=%s", status, w.Body.String())
	}
}

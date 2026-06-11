package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/thetramp22/rifflog/internal/models"
)

func TestPracticeSessions(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	app := SetupTestApp(t)
	defer app.DB.Close(context.Background())

	// Test User Setup
	email := fmt.Sprintf("test-%d@test.com", time.Now().UnixNano())
	password := "1234"

	user, err := CreateTestUser(app.UserRepo, email, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	t.Logf("registered user id=%d", user.ID)

	// Test 1 - valid request
	t.Log("creating request: valid request")
	practicedAt := time.Date(
		2026,
		time.June,
		10,
		14,
		30,
		0,
		0,
		time.UTC,
	)

	data := models.CreatePracticeSessionRequest{
		SkillID:         1,
		DurationMinutes: 20,
		PracticedAt:     practicedAt,
		Notes:           "short practice session",
		UserID:          user.ID,
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	req := httptest.NewRequest("POST", "http://localhost:8080/practice-sessions", bytes.NewReader(jsonBytes))
	w := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Fatalf("expected 201, got %v, body=%s", status, w.Body.String())
	}

	want := models.PracticeSession{
		SkillID:         1,
		DurationMinutes: 20,
		PracticedAt:     practicedAt,
		Notes:           "short practice session",
		UserID:          user.ID,
	}

	var got models.PracticeSession
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	opts := cmpopts.IgnoreFields(models.PracticeSession{}, "ID", "CreatedAt")
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("Values mismatch (-want +got):\n%s", diff)
	}

	// Test 2 - missing duration
	t.Log("creating request: missing duration")
	practicedAt = time.Date(
		2026,
		time.June,
		10,
		14,
		30,
		0,
		0,
		time.UTC,
	)

	data = models.CreatePracticeSessionRequest{
		SkillID:     1,
		PracticedAt: practicedAt,
		Notes:       "short practice session",
		UserID:      user.ID,
	}
	jsonBytes, err = json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	req = httptest.NewRequest("POST", "http://localhost:8080/practice-sessions", bytes.NewReader(jsonBytes))
	w = httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Fatalf("expected 400, got %v, body=%s", status, w.Body.String())
	}
}

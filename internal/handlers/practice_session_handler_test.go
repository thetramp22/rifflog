package handlers

import (
	"bytes"
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

func TestCreatePracticeSession(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	// Test 1 - Create Session
	t.Log("creating request: Create Session - Valid Request")
	practicedAt := time.Date(2026, time.June, 10, 14, 30, 0, 0, time.UTC)

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
}

func TestCreatePracticeSession_InvalidDuration(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	// Test 2 - Create Session - missing duration
	t.Log("creating request: Create Session - missing duration")
	practicedAt := time.Date(2026, time.June, 10, 14, 30, 0, 0, time.UTC)

	data := models.CreatePracticeSessionRequest{
		SkillID:     1,
		PracticedAt: practicedAt,
		Notes:       "short practice session",
		UserID:      user.ID,
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

	if status := w.Code; status != http.StatusBadRequest {
		t.Fatalf("expected 400, got %v, body=%s", status, w.Body.String())
	}
}

func TestListPracticeSessions(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	// Test Sessions Setup
	sessions := SeedPracticeSessionsForTest(t, app, user)
	practicedAt := sessions.PracticedAt
	scalesPracticedAt := sessions.ScalesPracticedAt
	laterEarTrainingPracticedAt := sessions.LaterEarTrainingPracticedAt

	// Test 3 - GET practice sessions
	want := []models.PracticeSessionDetails{
		{SkillID: 1,
			SkillName:        "Ear Training",
			SkillDescription: "Try playing to identify chords and melodies by ear.",
			DurationMinutes:  45,
			PracticedAt:      laterEarTrainingPracticedAt,
			Notes:            "long ear training session",
			UserID:           user.ID},
		{SkillID: 2,
			SkillName:        "Scales",
			SkillDescription: "Memorize note locations and scale patterns.",
			DurationMinutes:  35,
			PracticedAt:      scalesPracticedAt,
			Notes:            "scales practice",
			UserID:           user.ID},
		{SkillID: 1,
			SkillName:        "Ear Training",
			SkillDescription: "Try playing to identify chords and melodies by ear.",
			DurationMinutes:  20,
			PracticedAt:      practicedAt,
			Notes:            "short practice session",
			UserID:           user.ID},
	}

	got := getPracticeSessionsForTest(t, app, "http://localhost:8080/practice-sessions", user.ID)
	opts := cmpopts.IgnoreFields(models.PracticeSessionDetails{}, "ID", "CreatedAt")
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("Values mismatch (-want +got):\n%s", diff)
	}
}

func TestListPracticeSessions_FilterBySkill(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	// Test Sessions Setup
	sessions := SeedPracticeSessionsForTest(t, app, user)
	scalesPracticedAt := sessions.ScalesPracticedAt

	// Test 4 - GET practice sessions filtered by skill
	want := []models.PracticeSessionDetails{
		{SkillID: 2,
			SkillName:        "Scales",
			SkillDescription: "Memorize note locations and scale patterns.",
			DurationMinutes:  35,
			PracticedAt:      scalesPracticedAt,
			Notes:            "scales practice",
			UserID:           user.ID},
	}

	got := getPracticeSessionsQueryOnlyForTest(t, app, fmt.Sprintf("http://localhost:8080/practice-sessions?user_id=%d&skill=2", user.ID))
	opts := cmpopts.IgnoreFields(models.PracticeSessionDetails{}, "ID", "CreatedAt")
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("Values mismatch (-want +got):\n%s", diff)
	}
}

func TestListPracticeSessions_FilterByFromDate(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	// Test Sessions Setup
	sessions := SeedPracticeSessionsForTest(t, app, user)
	scalesPracticedAt := sessions.ScalesPracticedAt
	laterEarTrainingPracticedAt := sessions.LaterEarTrainingPracticedAt

	// Test 5 - GET practice sessions filtered by from date
	want := []models.PracticeSessionDetails{
		{SkillID: 1,
			SkillName:        "Ear Training",
			SkillDescription: "Try playing to identify chords and melodies by ear.",
			DurationMinutes:  45,
			PracticedAt:      laterEarTrainingPracticedAt,
			Notes:            "long ear training session",
			UserID:           user.ID},
		{SkillID: 2,
			SkillName:        "Scales",
			SkillDescription: "Memorize note locations and scale patterns.",
			DurationMinutes:  35,
			PracticedAt:      scalesPracticedAt,
			Notes:            "scales practice",
			UserID:           user.ID},
	}

	got := getPracticeSessionsForTest(t, app, "http://localhost:8080/practice-sessions?from=2026-06-11", user.ID)
	opts := cmpopts.IgnoreFields(models.PracticeSessionDetails{}, "ID", "CreatedAt")
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("Values mismatch (-want +got):\n%s", diff)
	}
}

func TestListPracticeSessions_FilterByToDate(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	// Test Sessions Setup
	sessions := SeedPracticeSessionsForTest(t, app, user)
	practicedAt := sessions.PracticedAt
	scalesPracticedAt := sessions.ScalesPracticedAt

	// Test 6 - GET practice sessions filtered by inclusive to date
	want := []models.PracticeSessionDetails{
		{SkillID: 2,
			SkillName:        "Scales",
			SkillDescription: "Memorize note locations and scale patterns.",
			DurationMinutes:  35,
			PracticedAt:      scalesPracticedAt,
			Notes:            "scales practice",
			UserID:           user.ID},
		{SkillID: 1,
			SkillName:        "Ear Training",
			SkillDescription: "Try playing to identify chords and melodies by ear.",
			DurationMinutes:  20,
			PracticedAt:      practicedAt,
			Notes:            "short practice session",
			UserID:           user.ID},
	}

	got := getPracticeSessionsForTest(t, app, "http://localhost:8080/practice-sessions?to=2026-06-11", user.ID)
	opts := cmpopts.IgnoreFields(models.PracticeSessionDetails{}, "ID", "CreatedAt")
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("Values mismatch (-want +got):\n%s", diff)
	}
}

func TestListPracticeSessionStats(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	// Test Sessions Setup
	sessions := SeedPracticeSessionsForTest(t, app, user)
	totalMinutes := sessions.TotalMinutes
	totalSessions := sessions.TotalSessions
	mostPracticedSkill := sessions.MostPracticedSkill
	longestSession := sessions.LongestSession

	want := models.PracticeSessionStats{
		TotalMinutes:       totalMinutes,
		TotalSessions:      totalSessions,
		MostPracticedSkill: mostPracticedSkill,
		LongestSession:     longestSession,
	}

	got := getPracticeSessionStatsForTest(t, app, "http://localhost:8080/practice-sessions/stats", user.ID)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Values mismatch (-want +got):\n%s", diff)
	}
}

func createPracticeSessionForTest(t *testing.T, app *TestApp, data models.CreatePracticeSessionRequest) {
	t.Helper()

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "http://localhost:8080/practice-sessions", bytes.NewReader(jsonBytes))
	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Fatalf("expected 201, got %v, body=%s", status, w.Body.String())
	}
}

func getPracticeSessionsForTest(t *testing.T, app *TestApp, target string, userID int) []models.PracticeSessionDetails {
	t.Helper()

	data := models.PracticeSessionDetailsRequest{
		UserID: userID,
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("GET", target, bytes.NewReader(jsonBytes))
	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("expected 200, got %v, body=%s", status, w.Body.String())
	}

	var got []models.PracticeSessionDetails
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	return got
}

func getPracticeSessionsQueryOnlyForTest(t *testing.T, app *TestApp, target string) []models.PracticeSessionDetails {
	t.Helper()

	req := httptest.NewRequest("GET", target, nil)
	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("expected 200, got %v, body=%s", status, w.Body.String())
	}

	var got []models.PracticeSessionDetails
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	return got
}

type SessionFixture struct {
	PracticedAt                 time.Time
	ScalesPracticedAt           time.Time
	LaterEarTrainingPracticedAt time.Time
	TotalMinutes                int
	TotalSessions               int
	MostPracticedSkill          models.MostPracticedSkill
	LongestSession              int
}

func SeedPracticeSessionsForTest(
	t *testing.T,
	app *TestApp,
	user models.User,
) SessionFixture {
	duration1 := 20
	duration2 := 35
	duration3 := 45
	totalMinutes := duration1 + duration2 + duration3

	practicedAt := time.Date(2026, time.June, 10, 14, 30, 0, 0, time.UTC)
	createPracticeSessionForTest(t, app, models.CreatePracticeSessionRequest{
		SkillID:         1,
		DurationMinutes: duration1,
		PracticedAt:     practicedAt,
		Notes:           "short practice session",
		UserID:          user.ID,
	})

	scalesPracticedAt := time.Date(2026, time.June, 11, 9, 0, 0, 0, time.UTC)
	createPracticeSessionForTest(t, app, models.CreatePracticeSessionRequest{
		SkillID:         2,
		DurationMinutes: duration2,
		PracticedAt:     scalesPracticedAt,
		Notes:           "scales practice",
		UserID:          user.ID,
	})

	laterEarTrainingPracticedAt := time.Date(2026, time.June, 12, 18, 45, 0, 0, time.UTC)
	createPracticeSessionForTest(t, app, models.CreatePracticeSessionRequest{
		SkillID:         1,
		DurationMinutes: duration3,
		PracticedAt:     laterEarTrainingPracticedAt,
		Notes:           "long ear training session",
		UserID:          user.ID,
	})

	mostPracticedSkill := models.MostPracticedSkill{
		Name:         "Ear Training",
		TotalMinutes: duration1 + duration3,
	}

	sessionFixture := SessionFixture{
		PracticedAt:                 practicedAt,
		ScalesPracticedAt:           scalesPracticedAt,
		LaterEarTrainingPracticedAt: laterEarTrainingPracticedAt,
		TotalMinutes:                totalMinutes,
		TotalSessions:               3,
		MostPracticedSkill:          mostPracticedSkill,
		LongestSession:              duration3,
	}

	return sessionFixture
}

func getPracticeSessionStatsForTest(t *testing.T, app *TestApp, target string, userID int) models.PracticeSessionStats {
	t.Helper()

	data := models.PracticeSessionStatsRequest{
		UserID: userID,
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("GET", target, bytes.NewReader(jsonBytes))
	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("expected 200, got %v, body=%s", status, w.Body.String())
	}

	var got models.PracticeSessionStats
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	return got
}

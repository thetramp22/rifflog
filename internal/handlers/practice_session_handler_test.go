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

	t.Log("creating request: Create Session - Valid Request")
	practicedAt := time.Date(2026, time.June, 10, 14, 30, 0, 0, time.UTC)

	data := models.CreatePracticeSessionRequest{
		SkillID:         1,
		DurationMinutes: 20,
		PracticedAt:     practicedAt,
		Notes:           "short practice session",
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	req := httptest.NewRequest("POST", "http://localhost:8080/api/practice-sessions", bytes.NewReader(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+user.Token)
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
		UserID:          user.User.ID,
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

	t.Log("creating request: Create Session - missing duration")
	practicedAt := time.Date(2026, time.June, 10, 14, 30, 0, 0, time.UTC)

	data := models.CreatePracticeSessionRequest{
		SkillID:     1,
		PracticedAt: practicedAt,
		Notes:       "short practice session",
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	req := httptest.NewRequest("POST", "http://localhost:8080/api/practice-sessions", bytes.NewReader(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+user.Token)
	w := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Fatalf("expected 400, got %v, body=%s", status, w.Body.String())
	}
}

func TestUpdatePracticeSession(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	t.Log("creating request: Create Session to be updated")
	practicedAt := time.Date(2026, time.June, 10, 14, 30, 0, 0, time.UTC)

	data := models.CreatePracticeSessionRequest{
		SkillID:         1,
		DurationMinutes: 20,
		PracticedAt:     practicedAt,
		Notes:           "short practice session",
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	req := httptest.NewRequest("POST", "http://localhost:8080/api/practice-sessions", bytes.NewReader(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+user.Token)
	w := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(w, req)

	var originalSession models.PracticeSession
	err = json.Unmarshal(w.Body.Bytes(), &originalSession)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	updateData := models.UpdatePracticeSessionRequest{
		SkillID:         2,
		DurationMinutes: 25,
		PracticedAt:     practicedAt,
		Notes:           "short practice session, edited skill id and duration",
	}
	jsonBytes, err = json.Marshal(updateData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	updateReq := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/api/practice-sessions/%v", originalSession.ID), bytes.NewReader(jsonBytes))
	updateReq.Header.Set("Authorization", "Bearer "+user.Token)
	updateW := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(updateW, updateReq)

	if status := updateW.Code; status != http.StatusOK {
		t.Fatalf("expected 200, got %v, body=%s", status, updateW.Body.String())
	}

	want := models.PracticeSession{
		SkillID:         2,
		DurationMinutes: 25,
		PracticedAt:     practicedAt,
		Notes:           "short practice session, edited skill id and duration",
		UserID:          user.User.ID,
	}

	var got models.PracticeSession
	err = json.Unmarshal(updateW.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	opts := cmpopts.IgnoreFields(models.PracticeSession{}, "ID", "CreatedAt")
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("Values mismatch (-want +got):\n%s", diff)
	}
}

func TestUpdatePracticeSessionWrongUser(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	wrongUser := SetupExtraTestUser(t, app, "test2")
	defer app.DB.Close()

	t.Log("creating request: Create Session to be updated")
	practicedAt := time.Date(2026, time.June, 10, 14, 30, 0, 0, time.UTC)

	data := models.CreatePracticeSessionRequest{
		SkillID:         1,
		DurationMinutes: 20,
		PracticedAt:     practicedAt,
		Notes:           "short practice session",
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	req := httptest.NewRequest("POST", "http://localhost:8080/api/practice-sessions", bytes.NewReader(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+user.Token)
	w := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(w, req)

	var originalSession models.PracticeSession
	err = json.Unmarshal(w.Body.Bytes(), &originalSession)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	t.Log("creating request: Updated Session")
	updateData := models.UpdatePracticeSessionRequest{
		SkillID:         2,
		DurationMinutes: 25,
		PracticedAt:     practicedAt,
		Notes:           "short practice session, edited skill id and duration",
	}
	jsonBytes, err = json.Marshal(updateData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	updateReq := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/api/practice-sessions/%v", originalSession.ID), bytes.NewReader(jsonBytes))
	updateReq.Header.Set("Authorization", "Bearer "+wrongUser.Token)
	updateW := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(updateW, updateReq)

	if status := updateW.Code; status != http.StatusNotFound {
		t.Fatalf("expected 404, got %v, body=%s", status, updateW.Body.String())
	}
}

func TestDeletePracticeSession(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	t.Log("creating request: Create Session to be deleted")
	practicedAt := time.Date(2026, time.June, 10, 14, 30, 0, 0, time.UTC)

	data := models.CreatePracticeSessionRequest{
		SkillID:         1,
		DurationMinutes: 20,
		PracticedAt:     practicedAt,
		Notes:           "short practice session",
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	req := httptest.NewRequest("POST", "http://localhost:8080/api/practice-sessions", bytes.NewReader(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+user.Token)
	w := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(w, req)

	var originalSession models.PracticeSession
	err = json.Unmarshal(w.Body.Bytes(), &originalSession)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	t.Log(originalSession)

	t.Log("creating request: Delete Session")
	deleteReq := httptest.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/api/practice-sessions/%v", originalSession.ID), nil)
	deleteReq.Header.Set("Authorization", "Bearer "+user.Token)
	deleteW := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(deleteW, deleteReq)

	if status := deleteW.Code; status != http.StatusOK {
		t.Fatalf("expected 200, got %v, body=%s", status, deleteW.Body.String())
	}
}

func TestDeletePracticeSessionWrongUser(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	wrongUser := SetupExtraTestUser(t, app, "test2")
	defer app.DB.Close()

	t.Log("creating request: Create Session to be deleted")
	practicedAt := time.Date(2026, time.June, 10, 14, 30, 0, 0, time.UTC)

	data := models.CreatePracticeSessionRequest{
		SkillID:         1,
		DurationMinutes: 20,
		PracticedAt:     practicedAt,
		Notes:           "short practice session",
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	t.Log(string(jsonBytes))
	req := httptest.NewRequest("POST", "http://localhost:8080/api/practice-sessions", bytes.NewReader(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+user.Token)
	w := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(w, req)

	var originalSession models.PracticeSession
	err = json.Unmarshal(w.Body.Bytes(), &originalSession)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	t.Log(originalSession)

	t.Log("creating request: Delete Session")
	deleteReq := httptest.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/api/practice-sessions/%v", originalSession.ID), nil)
	deleteReq.Header.Set("Authorization", "Bearer "+wrongUser.Token)
	deleteW := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(deleteW, deleteReq)

	if status := deleteW.Code; status != http.StatusNotFound {
		t.Fatalf("expected 404, got %v, body=%s", status, deleteW.Body.String())
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

	want := []models.PracticeSessionDetails{
		{SkillID: 1,
			SkillName:        "Ear Training",
			SkillDescription: "Try playing to identify chords and melodies by ear.",
			DurationMinutes:  45,
			PracticedAt:      laterEarTrainingPracticedAt,
			Notes:            "long ear training session",
			UserID:           user.User.ID},
		{SkillID: 2,
			SkillName:        "Scales",
			SkillDescription: "Memorize note locations and scale patterns.",
			DurationMinutes:  35,
			PracticedAt:      scalesPracticedAt,
			Notes:            "scales practice",
			UserID:           user.User.ID},
		{SkillID: 1,
			SkillName:        "Ear Training",
			SkillDescription: "Try playing to identify chords and melodies by ear.",
			DurationMinutes:  20,
			PracticedAt:      practicedAt,
			Notes:            "short practice session",
			UserID:           user.User.ID},
	}

	got := getPracticeSessionsForTest(t, app, user, "http://localhost:8080/api/practice-sessions")
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

	want := []models.PracticeSessionDetails{
		{SkillID: 2,
			SkillName:        "Scales",
			SkillDescription: "Memorize note locations and scale patterns.",
			DurationMinutes:  35,
			PracticedAt:      scalesPracticedAt,
			Notes:            "scales practice",
			UserID:           user.User.ID},
	}

	got := getPracticeSessionsForTest(t, app, user, "http://localhost:8080/api/practice-sessions?skill=2")
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

	want := []models.PracticeSessionDetails{
		{SkillID: 1,
			SkillName:        "Ear Training",
			SkillDescription: "Try playing to identify chords and melodies by ear.",
			DurationMinutes:  45,
			PracticedAt:      laterEarTrainingPracticedAt,
			Notes:            "long ear training session",
			UserID:           user.User.ID},
		{SkillID: 2,
			SkillName:        "Scales",
			SkillDescription: "Memorize note locations and scale patterns.",
			DurationMinutes:  35,
			PracticedAt:      scalesPracticedAt,
			Notes:            "scales practice",
			UserID:           user.User.ID},
	}

	got := getPracticeSessionsForTest(t, app, user, "http://localhost:8080/api/practice-sessions?from=2026-06-11")
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

	want := []models.PracticeSessionDetails{
		{SkillID: 2,
			SkillName:        "Scales",
			SkillDescription: "Memorize note locations and scale patterns.",
			DurationMinutes:  35,
			PracticedAt:      scalesPracticedAt,
			Notes:            "scales practice",
			UserID:           user.User.ID},
		{SkillID: 1,
			SkillName:        "Ear Training",
			SkillDescription: "Try playing to identify chords and melodies by ear.",
			DurationMinutes:  20,
			PracticedAt:      practicedAt,
			Notes:            "short practice session",
			UserID:           user.User.ID},
	}

	got := getPracticeSessionsForTest(t, app, user, "http://localhost:8080/api/practice-sessions?to=2026-06-11")
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
		MostPracticedSkill: &mostPracticedSkill,
		LongestSession:     longestSession,
	}

	got := getPracticeSessionStatsForTest(t, app, user, "http://localhost:8080/api/practice-sessions/stats")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Values mismatch (-want +got):\n%s", diff)
	}
}

func TestListPracticeSessionsNoSessions(t *testing.T) {
	// Test App Setup
	t.Log("creating router")
	password := "test"
	app, user := SetupTestUser(t, password)
	defer app.DB.Close()

	want := models.PracticeSessionStats{
		TotalMinutes:       0,
		TotalSessions:      0,
		MostPracticedSkill: nil,
		LongestSession:     0,
	}

	got := getPracticeSessionStatsForTest(t, app, user, "http://localhost:8080/api/practice-sessions/stats")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Values mismatch (-want +got):\n%s", diff)
	}
}

func createPracticeSessionForTest(t *testing.T, app *TestApp, user TestUser, data models.CreatePracticeSessionRequest) {
	t.Helper()

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "http://localhost:8080/api/practice-sessions", bytes.NewReader(jsonBytes))
	w := httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer "+user.Token)
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Fatalf("expected 201, got %v, body=%s", status, w.Body.String())
	}
}

func getPracticeSessionsForTest(t *testing.T, app *TestApp, user TestUser, target string) []models.PracticeSessionDetails {
	t.Helper()

	req := httptest.NewRequest("GET", target, nil)
	req.Header.Set("Authorization", "Bearer "+user.Token)
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
	user TestUser,
) SessionFixture {
	duration1 := 20
	duration2 := 35
	duration3 := 45
	totalMinutes := duration1 + duration2 + duration3

	practicedAt := time.Date(2026, time.June, 10, 14, 30, 0, 0, time.UTC)
	createPracticeSessionForTest(t, app, user, models.CreatePracticeSessionRequest{
		SkillID:         1,
		DurationMinutes: duration1,
		PracticedAt:     practicedAt,
		Notes:           "short practice session",
	})

	scalesPracticedAt := time.Date(2026, time.June, 11, 9, 0, 0, 0, time.UTC)
	createPracticeSessionForTest(t, app, user, models.CreatePracticeSessionRequest{
		SkillID:         2,
		DurationMinutes: duration2,
		PracticedAt:     scalesPracticedAt,
		Notes:           "scales practice",
	})

	laterEarTrainingPracticedAt := time.Date(2026, time.June, 12, 18, 45, 0, 0, time.UTC)
	createPracticeSessionForTest(t, app, user, models.CreatePracticeSessionRequest{
		SkillID:         1,
		DurationMinutes: duration3,
		PracticedAt:     laterEarTrainingPracticedAt,
		Notes:           "long ear training session",
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

func getPracticeSessionStatsForTest(t *testing.T, app *TestApp, user TestUser, target string) models.PracticeSessionStats {
	t.Helper()

	req := httptest.NewRequest("GET", target, nil)
	req.Header.Set("Authorization", "Bearer "+user.Token)
	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("expected 200, got %v, body=%s", status, w.Body.String())
	}

	var got models.PracticeSessionStats
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	return got
}

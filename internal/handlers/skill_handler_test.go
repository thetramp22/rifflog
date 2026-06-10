package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/thetramp22/rifflog/internal/models"
)

func TestSkillsEndpoint(t *testing.T) {
	t.Log("creating router")
	app := SetupTestApp(t)
	defer app.DB.Close(context.Background())

	t.Log("creating request")
	req := httptest.NewRequest("GET", "http://localhost:8080/skills", nil)
	w := httptest.NewRecorder()

	t.Log("ServeHTTP call")
	app.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected 200, got %v", status)
	}

	want := []models.Skill{
		{
			Name:        "Ear Training",
			Description: "Try playing to identify chords and melodies by ear.",
		},
		{
			Name:        "Scales",
			Description: "Memorize note locations and scale patterns.",
		},
		{
			Name:        "Timing and Rhythm",
			Description: "Practice with a metronome to develop a solid sense of time and groove.",
		},
	}

	var got []models.Skill
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	opts := cmpopts.IgnoreFields(models.Skill{}, "ID", "CreatedAt")
	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("Values mismatch (-want +got):\n%s", diff)
	}
}

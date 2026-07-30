package handlers

import (
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
	defer app.DB.Close()

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
			Name:        "Alternate Picking",
			Description: "Alternate Picking is a picking technique that alternates downstrokes and upstrokes to improve speed, rhythm, and efficiency.",
		},
		{
			Name:        "Chord Transition",
			Description: "Chord Transition is moving your fretting hand smoothly and efficiently from one chord shape to another in time with the music.",
		},
		{
			Name:        "Ear Training",
			Description: "Ear training is the process of learning to recognize and identify musical pitches, intervals, chords, and rhythms by ear.",
		},
		{
			Name:        "Finger Independence",
			Description: "Finger independence is the ability to control and move each finger on your fretting (or picking) hand completely isolated from the others.",
		},
		{
			Name:        "Scales",
			Description: "A guitar scale is asequence of notes played in ascending or descending order of pitch, separated by specific intervals of whole and half steps.",
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

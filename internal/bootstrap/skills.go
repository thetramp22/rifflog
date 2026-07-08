// Package bootstrap provides setup for the application on startup
package bootstrap

import (
	"context"

	"github.com/thetramp22/rifflog/internal/models"
	"github.com/thetramp22/rifflog/internal/repository"
)

// PopulateSkillsList creates the list of skills used by the app.
// These skills are used to create practice sessions by the user.
// In future iteration this function may be replaced by a feature to add
// and/or remove custom skills by the user.
func PopulateSkillsList(ctx context.Context, r *repository.SkillRepository) error {
	skills := []models.Skill{
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

	for _, skill := range skills {
		err := r.SeedSkill(ctx, skill)
		if err != nil {
			return err
		}
	}

	return nil
}

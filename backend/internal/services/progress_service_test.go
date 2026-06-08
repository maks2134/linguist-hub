package services

import (
	"testing"

	"github.com/makskozlov/linguist-hub/internal/models"
)

func TestCountPlanProgress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		modules       models.Modules
		wantCompleted int
		wantTotal     int
	}{
		{
			name:          "empty modules",
			modules:       models.Modules{},
			wantCompleted: 0,
			wantTotal:     0,
		},
		{
			name: "module with no topics counts zero total",
			modules: models.Modules{
				{ID: 1, Topics: []string{}, CompletedTopics: []string{"orphan"}},
			},
			wantCompleted: 0,
			wantTotal:     0,
		},
		{
			name: "partial completion across modules",
			modules: models.Modules{
				{ID: 1, Topics: []string{"A", "B"}, CompletedTopics: []string{"A"}},
				{ID: 2, Topics: []string{"C"}, CompletedTopics: []string{"C"}},
			},
			wantCompleted: 2,
			wantTotal:     3,
		},
		{
			name: "whitespace normalized when counting",
			modules: models.Modules{
				{ID: 1, Topics: []string{" A ", "B"}, CompletedTopics: []string{"A", " B "}},
			},
			wantCompleted: 2,
			wantTotal:     2,
		},
		{
			name: "completed topic not in plan topics ignored",
			modules: models.Modules{
				{ID: 1, Topics: []string{"A"}, CompletedTopics: []string{"A", "extra"}},
			},
			wantCompleted: 1,
			wantTotal:     1,
		},
		{
			name: "multiple modules aggregate totals",
			modules: models.Modules{
				{ID: 1, Topics: []string{"A", "B"}, CompletedTopics: []string{"A", "B"}},
				{ID: 2, Topics: []string{"C", "D", "E"}, CompletedTopics: []string{"C"}},
			},
			wantCompleted: 3,
			wantTotal:     5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotCompleted, gotTotal := CountPlanProgress(tt.modules)
			if gotCompleted != tt.wantCompleted || gotTotal != tt.wantTotal {
				t.Fatalf("CountPlanProgress() = (%d, %d), want (%d, %d)",
					gotCompleted, gotTotal, tt.wantCompleted, tt.wantTotal)
			}
		})
	}
}

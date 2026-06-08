package services

import (
	"testing"

	"github.com/makskozlov/linguist-hub/internal/models"
)

func TestMarkTopicCompleted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		modules    models.Modules
		moduleID   int
		topic      string
		wantMarked bool
		wantCount  int
	}{
		{
			name: "new topic marked in matching module",
			modules: models.Modules{
				{ID: 1, Topics: []string{"A", "B"}, CompletedTopics: []string{}},
			},
			moduleID:   1,
			topic:      "A",
			wantMarked: true,
			wantCount:  1,
		},
		{
			name: "duplicate topic not marked again",
			modules: models.Modules{
				{ID: 1, Topics: []string{"A", "B"}, CompletedTopics: []string{"A"}},
			},
			moduleID:   1,
			topic:      "A",
			wantMarked: false,
			wantCount:  1,
		},
		{
			name: "whitespace normalized for duplicate detection",
			modules: models.Modules{
				{ID: 1, Topics: []string{"A"}, CompletedTopics: []string{" A "}},
			},
			moduleID:   1,
			topic:      "A",
			wantMarked: false,
			wantCount:  1,
		},
		{
			name: "empty topic rejected",
			modules: models.Modules{
				{ID: 1, Topics: []string{"A"}},
			},
			moduleID:   1,
			topic:      "  ",
			wantMarked: false,
			wantCount:  0,
		},
		{
			name: "unknown module id",
			modules: models.Modules{
				{ID: 1, Topics: []string{"A"}},
			},
			moduleID:   99,
			topic:      "A",
			wantMarked: false,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			modules := append(models.Modules(nil), tt.modules...)
			got := MarkTopicCompleted(&modules, tt.moduleID, tt.topic)
			if got != tt.wantMarked {
				t.Fatalf("MarkTopicCompleted() = %v, want %v", got, tt.wantMarked)
			}
			if tt.moduleID == 1 && len(modules[0].CompletedTopics) != tt.wantCount {
				t.Fatalf("completed topics count = %d, want %d", len(modules[0].CompletedTopics), tt.wantCount)
			}
		})
	}
}

func TestIsModuleComplete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		module models.Module
		want   bool
	}{
		{
			name:   "empty topics is complete",
			module: models.Module{Topics: []string{}},
			want:   true,
		},
		{
			name:   "all topics completed",
			module: models.Module{Topics: []string{"A", "B"}, CompletedTopics: []string{"A", "B"}},
			want:   true,
		},
		{
			name:   "partial completion",
			module: models.Module{Topics: []string{"A", "B"}, CompletedTopics: []string{"A"}},
			want:   false,
		},
		{
			name:   "whitespace normalized",
			module: models.Module{Topics: []string{"A"}, CompletedTopics: []string{" A "}},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := IsModuleComplete(tt.module); got != tt.want {
				t.Fatalf("IsModuleComplete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleLocked(t *testing.T) {
	t.Parallel()

	modules := models.Modules{
		{ID: 1, Topics: []string{"A"}, CompletedTopics: []string{"A"}},
		{ID: 2, Topics: []string{"B"}, CompletedTopics: []string{}},
		{ID: 3, Topics: []string{"C"}, CompletedTopics: []string{}},
	}

	tests := []struct {
		name  string
		index int
		want  bool
	}{
		{name: "first module unlocked", index: 0, want: false},
		{name: "second module unlocked when first complete", index: 1, want: false},
		{name: "third module locked when second incomplete", index: 2, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ModuleLocked(modules, tt.index); got != tt.want {
				t.Fatalf("ModuleLocked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsModuleLocked(t *testing.T) {
	t.Parallel()

	modules := models.Modules{
		{ID: 10, Topics: []string{"A"}, CompletedTopics: []string{"A"}},
		{ID: 20, Topics: []string{"B"}, CompletedTopics: []string{}},
	}

	tests := []struct {
		name     string
		moduleID int
		want     bool
	}{
		{name: "first module by id unlocked", moduleID: 10, want: false},
		{name: "second module by id unlocked", moduleID: 20, want: false},
		{name: "unknown module id treated as locked", moduleID: 999, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := IsModuleLocked(modules, tt.moduleID); got != tt.want {
				t.Fatalf("IsModuleLocked() = %v, want %v", got, tt.want)
			}
		})
	}
}

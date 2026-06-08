package repository

import (
	"testing"
	"time"

	"github.com/makskozlov/linguist-hub/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupVoiceLessonTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.VoiceLesson{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestVoiceLessonRepository_MarkAbandonedStale(t *testing.T) {
	now := time.Now()
	stale := now.Add(-2 * time.Hour)
	recent := now.Add(-5 * time.Minute)

	tests := []struct {
		name       string
		lessons    []models.VoiceLesson
		olderThan  time.Duration
		wantRows   int64
		wantStatus map[uint]string
	}{
		{
			name: "marks stale in-progress lessons",
			lessons: []models.VoiceLesson{
				{UserID: 1, ProfileID: 1, Topic: "t1", Status: models.VoiceLessonStatusInProgress, StartedAt: stale, LastTurnAt: stale},
				{UserID: 1, ProfileID: 1, Topic: "t2", Status: models.VoiceLessonStatusInProgress, StartedAt: recent, LastTurnAt: recent},
			},
			olderThan: time.Hour,
			wantRows:  1,
			wantStatus: map[uint]string{
				1: models.VoiceLessonStatusAbandoned,
				2: models.VoiceLessonStatusInProgress,
			},
		},
		{
			name: "does not touch completed lessons",
			lessons: []models.VoiceLesson{
				{UserID: 1, ProfileID: 1, Topic: "t1", Status: models.VoiceLessonStatusCompleted, StartedAt: stale, LastTurnAt: stale},
			},
			olderThan: time.Hour,
			wantRows:  0,
			wantStatus: map[uint]string{
				1: models.VoiceLessonStatusCompleted,
			},
		},
		{
			name: "does not re-mark already abandoned",
			lessons: []models.VoiceLesson{
				{UserID: 1, ProfileID: 1, Topic: "t1", Status: models.VoiceLessonStatusAbandoned, StartedAt: stale, LastTurnAt: stale},
			},
			olderThan: time.Hour,
			wantRows:  0,
			wantStatus: map[uint]string{
				1: models.VoiceLessonStatusAbandoned,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupVoiceLessonTestDB(t)
			repo := NewVoiceLessonRepository(db)

			for i := range tt.lessons {
				lesson := tt.lessons[i]
				lesson.TargetPhrases = models.JSONStringList{}
				lesson.HitPhrases = models.JSONStringList{}
				if err := repo.Create(&lesson); err != nil {
					t.Fatalf("create lesson: %v", err)
				}
				tt.lessons[i].ID = lesson.ID
			}

			gotRows, err := repo.MarkAbandonedStale(tt.olderThan)
			if err != nil {
				t.Fatalf("MarkAbandonedStale() error = %v", err)
			}
			if gotRows != tt.wantRows {
				t.Fatalf("MarkAbandonedStale() rows = %d, want %d", gotRows, tt.wantRows)
			}

			for id, wantStatus := range tt.wantStatus {
				lesson, err := repo.FindByID(id)
				if err != nil {
					t.Fatalf("FindByID(%d): %v", id, err)
				}
				if lesson.Status != wantStatus {
					t.Fatalf("lesson %d status = %q, want %q", id, lesson.Status, wantStatus)
				}
			}
		})
	}
}

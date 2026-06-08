package repository

import (
	"time"

	"github.com/makskozlov/linguist-hub/internal/models"
	"gorm.io/gorm"
)

type VoiceLessonRepository struct {
	db *gorm.DB
}

func NewVoiceLessonRepository(db *gorm.DB) *VoiceLessonRepository {
	return &VoiceLessonRepository{db: db}
}

func (r *VoiceLessonRepository) Create(lesson *models.VoiceLesson) error {
	return r.db.Create(lesson).Error
}

func (r *VoiceLessonRepository) FindByID(id uint) (*models.VoiceLesson, error) {
	var lesson models.VoiceLesson
	err := r.db.First(&lesson, id).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (r *VoiceLessonRepository) FindActiveByProfile(profileID uint) (*models.VoiceLesson, error) {
	var lesson models.VoiceLesson
	err := r.db.Where("profile_id = ? AND status = ?", profileID, models.VoiceLessonStatusInProgress).
		Order("started_at DESC").
		First(&lesson).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (r *VoiceLessonRepository) ListByProfile(profileID uint, limit, offset int) ([]models.VoiceLesson, error) {
	var lessons []models.VoiceLesson
	q := r.db.Where("profile_id = ?", profileID).Order("started_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	err := q.Find(&lessons).Error
	return lessons, err
}

func (r *VoiceLessonRepository) Update(lesson *models.VoiceLesson) error {
	return r.db.Save(lesson).Error
}

func (r *VoiceLessonRepository) MarkAbandonedStale(olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	res := r.db.Model(&models.VoiceLesson{}).
		Where("status = ? AND last_turn_at < ?", models.VoiceLessonStatusInProgress, cutoff).
		Update("status", models.VoiceLessonStatusAbandoned)
	return res.RowsAffected, res.Error
}

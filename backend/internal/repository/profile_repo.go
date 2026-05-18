package repository

import (
	"github.com/makskozlov/linguist-hub/internal/models"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Create(profile *models.LanguageProfile) error {
	return r.db.Create(profile).Error
}

func (r *ProfileRepository) FindByUserIDAndLanguage(userID uint, language string) (*models.LanguageProfile, error) {
	var profile models.LanguageProfile
	err := r.db.Where("user_id = ? AND target_language = ?", userID, language).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) FindByID(id uint) (*models.LanguageProfile, error) {
	var profile models.LanguageProfile
	err := r.db.First(&profile, id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) FindByUserID(userID uint) ([]models.LanguageProfile, error) {
	var profiles []models.LanguageProfile
	err := r.db.Where("user_id = ?", userID).Find(&profiles).Error
	return profiles, err
}

func (r *ProfileRepository) Update(profile *models.LanguageProfile) error {
	return r.db.Save(profile).Error
}

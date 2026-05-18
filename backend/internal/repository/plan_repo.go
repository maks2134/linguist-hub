package repository

import (
	"github.com/makskozlov/linguist-hub/internal/models"
	"gorm.io/gorm"
)

type PlanRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) Create(plan *models.StudyPlan) error {
	return r.db.Create(plan).Error
}

func (r *PlanRepository) FindByProfileID(profileID uint) (*models.StudyPlan, error) {
	var plan models.StudyPlan
	err := r.db.Where("profile_id = ?", profileID).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) Update(plan *models.StudyPlan) error {
	return r.db.Save(plan).Error
}

package services

import (
	"errors"

	"github.com/makskozlov/linguist-hub/internal/constants"
	"github.com/makskozlov/linguist-hub/internal/models"
	"github.com/makskozlov/linguist-hub/internal/repository"
	"github.com/makskozlov/linguist-hub/pkg/apierror"
	"gorm.io/gorm"
)

type ProgressService struct {
	profileRepo *repository.ProfileRepository
	planRepo    *repository.PlanRepository
}

func NewProgressService(profileRepo *repository.ProfileRepository, planRepo *repository.PlanRepository) *ProgressService {
	return &ProgressService{
		profileRepo: profileRepo,
		planRepo:    planRepo,
	}
}

func CountPlanProgress(modules models.Modules) (completed, total int) {
	for _, m := range modules {
		total += len(m.Topics)
		completedSet := make(map[string]bool, len(m.CompletedTopics))
		for _, ct := range m.CompletedTopics {
			completedSet[normalizeTopic(ct)] = true
		}
		for _, t := range m.Topics {
			if completedSet[normalizeTopic(t)] {
				completed++
			}
		}
	}
	return completed, total
}

func (s *ProgressService) Recalculate(profileID uint) error {
	plan, err := s.planRepo.FindByProfileID(profileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	completed, total := CountPlanProgress(plan.Structure)
	var pct float64
	if total > 0 {
		pct = float64(completed) / float64(total) * 100
	}

	profile, err := s.profileRepo.FindByID(profileID)
	if err != nil {
		return apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	profile.ProgressPercentage = pct
	return s.profileRepo.Update(profile)
}

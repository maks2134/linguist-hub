package services

import (
	"github.com/makskozlov/linguist-hub/internal/dto"
	"github.com/makskozlov/linguist-hub/internal/models"
	"github.com/makskozlov/linguist-hub/internal/repository"
	"github.com/makskozlov/linguist-hub/pkg/apierror"
	"gorm.io/gorm"
)

type ProfileService struct {
	profileRepo *repository.ProfileRepository
}

func NewProfileService(profileRepo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{profileRepo: profileRepo}
}

func (s *ProfileService) SelectLanguage(userID uint, language string) (*dto.ProfileResponse, error) {
	if language == "" {
		return nil, apierror.New(400, "VALIDATION_ERROR", "Language is required")
	}

	existingProfile, err := s.profileRepo.FindByUserIDAndLanguage(userID, language)
	if err == nil && existingProfile != nil {
		return &dto.ProfileResponse{
			ID:                 existingProfile.ID,
			UserID:             existingProfile.UserID,
			TargetLanguage:     existingProfile.TargetLanguage,
			CurrentLevel:       existingProfile.CurrentLevel,
			ProgressPercentage: existingProfile.ProgressPercentage,
		}, nil
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, apierror.Wrap(err, 500, "DB_ERROR")
	}

	profile := &models.LanguageProfile{
		UserID:             userID,
		TargetLanguage:     language,
		CurrentLevel:       "unknown",
		ProgressPercentage: 0,
	}

	if err := s.profileRepo.Create(profile); err != nil {
		return nil, apierror.Wrap(err, 500, "DB_ERROR")
	}

	return &dto.ProfileResponse{
		ID:                 profile.ID,
		UserID:             profile.UserID,
		TargetLanguage:     profile.TargetLanguage,
		CurrentLevel:       profile.CurrentLevel,
		ProgressPercentage: profile.ProgressPercentage,
	}, nil
}

func (s *ProfileService) GetProfiles(userID uint) ([]dto.ProfileResponse, error) {
	profiles, err := s.profileRepo.FindByUserID(userID)
	if err != nil {
		return nil, apierror.Wrap(err, 500, "DB_ERROR")
	}

	result := make([]dto.ProfileResponse, len(profiles))
	for i, p := range profiles {
		result[i] = dto.ProfileResponse{
			ID:                 p.ID,
			UserID:             p.UserID,
			TargetLanguage:     p.TargetLanguage,
			CurrentLevel:       p.CurrentLevel,
			ProgressPercentage: p.ProgressPercentage,
		}
	}

	return result, nil
}

func (s *ProfileService) GetProfile(profileID uint) (*dto.ProfileResponse, error) {
	profile, err := s.profileRepo.FindByID(profileID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierror.ErrNotFound
		}
		return nil, apierror.Wrap(err, 500, "DB_ERROR")
	}

	return &dto.ProfileResponse{
		ID:                 profile.ID,
		UserID:             profile.UserID,
		TargetLanguage:     profile.TargetLanguage,
		CurrentLevel:       profile.CurrentLevel,
		ProgressPercentage: profile.ProgressPercentage,
	}, nil
}

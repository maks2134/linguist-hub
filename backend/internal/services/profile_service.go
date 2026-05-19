package services

import (
	"errors"

	"github.com/makskozlov/linguist-hub/internal/constants"
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

func profileToDTO(p *models.LanguageProfile) *dto.ProfileResponse {
	return &dto.ProfileResponse{
		ID:                 p.ID,
		UserID:             p.UserID,
		TargetLanguage:     p.TargetLanguage,
		NativeLanguage:     p.NativeLanguage,
		CurrentLevel:       p.CurrentLevel,
		OnboardingPhase:    p.OnboardingPhase,
		Goals:              p.Goals,
		ProgressPercentage: p.ProgressPercentage,
	}
}

func (s *ProfileService) SelectLanguage(userID uint, language, nativeLanguage string) (*dto.ProfileResponse, error) {
	if language == "" {
		return nil, apierror.ErrBadRequest
	}
	if nativeLanguage == "" {
		nativeLanguage = "English"
	}

	existingProfile, err := s.profileRepo.FindByUserIDAndLanguage(userID, language)
	if err == nil {
		return profileToDTO(existingProfile), nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	profile := &models.LanguageProfile{
		UserID:             userID,
		TargetLanguage:     language,
		NativeLanguage:     nativeLanguage,
		CurrentLevel:       "unknown",
		OnboardingPhase:    "goals",
		ProgressPercentage: 0,
	}

	if err := s.profileRepo.Create(profile); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	return profileToDTO(profile), nil
}

func (s *ProfileService) GetProfiles(userID uint) ([]dto.ProfileResponse, error) {
	profiles, err := s.profileRepo.FindByUserID(userID)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	result := make([]dto.ProfileResponse, len(profiles))
	for i, p := range profiles {
		result[i] = *profileToDTO(&p)
	}

	return result, nil
}

func (s *ProfileService) UpdateNativeLanguage(userID, profileID uint, nativeLanguage string) (*dto.ProfileResponse, error) {
	if nativeLanguage == "" {
		return nil, apierror.ErrBadRequest
	}

	profile, err := s.profileRepo.FindByID(profileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrNotFound
		}
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	if profile.UserID != userID {
		return nil, apierror.ErrForbidden
	}

	profile.NativeLanguage = nativeLanguage
	if err := s.profileRepo.Update(profile); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	return profileToDTO(profile), nil
}

func (s *ProfileService) GetProfile(userID, profileID uint) (*dto.ProfileResponse, error) {
	profile, err := s.profileRepo.FindByID(profileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrNotFound
		}
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	if profile.UserID != userID {
		return nil, apierror.ErrForbidden
	}

	return profileToDTO(profile), nil
}

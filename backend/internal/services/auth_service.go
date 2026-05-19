package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/makskozlov/linguist-hub/internal/constants"
	"github.com/makskozlov/linguist-hub/internal/dto"
	"github.com/makskozlov/linguist-hub/internal/models"
	"github.com/makskozlov/linguist-hub/internal/repository"
	"github.com/makskozlov/linguist-hub/pkg/apierror"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(email, password string) (*dto.AuthResponse, error) {
	if email == "" || password == "" {
		return nil, apierror.ErrBadRequest
	}

	if len(password) < constants.MinPasswordLength {
		return nil, apierror.New(constants.StatusBadRequest, "VALIDATION_ERROR", "Password must be at least 6 characters")
	}

	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return nil, apierror.ErrEmailAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "HASH_ERROR")
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "TOKEN_ERROR")
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.User{
			ID:       user.ID,
			Email:    user.Email,
			Nickname: user.Nickname,
		},
	}, nil
}

func (s *AuthService) Login(email, password string) (*dto.AuthResponse, error) {
	if email == "" || password == "" {
		return nil, apierror.ErrBadRequest
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrInvalidCredentials
		}
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, apierror.ErrInvalidCredentials
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "TOKEN_ERROR")
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.User{
			ID:       user.ID,
			Email:    user.Email,
			Nickname: user.Nickname,
		},
	}, nil
}

func (s *AuthService) UpdateProfile(userID uint, nickname string) (*dto.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrNotFound
		}
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	user.Nickname = nickname
	if err := s.userRepo.Update(user); err != nil {
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	return &dto.User{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
	}, nil
}

func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	if len(newPassword) < constants.MinPasswordLength {
		return apierror.New(constants.StatusBadRequest, "VALIDATION_ERROR", "Password must be at least 6 characters")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierror.ErrNotFound
		}
		return apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return apierror.ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apierror.Wrap(err, constants.StatusInternalError, "HASH_ERROR")
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(user)
}

func (s *AuthService) GetMe(userID uint) (*dto.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrNotFound
		}
		return nil, apierror.Wrap(err, constants.StatusInternalError, "DB_ERROR")
	}

	return &dto.User{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
	}, nil
}

func (s *AuthService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(constants.JWTTokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

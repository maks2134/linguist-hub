package repository

import (
	"github.com/makskozlov/linguist-hub/internal/models"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(message *models.ChatMessage) error {
	return r.db.Create(message).Error
}

func (r *ChatRepository) FindByUserIDAndContext(userID uint, context string, profileID uint) ([]models.ChatMessage, error) {
	return r.FindByUserIDAndContextPaginated(userID, context, profileID, 0, 0)
}

func (r *ChatRepository) FindByUserIDAndContextPaginated(userID uint, context string, profileID uint, limit, offset int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	query := r.db.Where("user_id = ? AND context = ?", userID, context)

	if profileID > 0 {
		query = query.Where("profile_id = ?", profileID)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at ASC").Find(&messages).Error
	return messages, err
}

func (r *ChatRepository) DeleteByUserIDAndContext(userID uint, context string, profileID uint) error {
	query := r.db.Where("user_id = ? AND context = ?", userID, context)

	if profileID > 0 {
		query = query.Where("profile_id = ?", profileID)
	}

	return query.Delete(&models.ChatMessage{}).Error
}

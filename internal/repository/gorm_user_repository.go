package repository

import (
	"errors"

	"ktx-diet/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) GetOrCreateByTelegramID(telegramID int64) (*domain.User, error) {
	user := &domain.User{TelegramID: telegramID}

	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "telegram_id"}},
		DoNothing: true,
	}).Create(user).Error
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		if err := r.db.Where("telegram_id = ?", telegramID).First(user).Error; err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, errors.New("failed to load user after get-or-create")
	}

	return user, nil
}

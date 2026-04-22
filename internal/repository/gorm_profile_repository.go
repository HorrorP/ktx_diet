package repository

import (
	"ktx-diet/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormProfileRepository struct {
	db *gorm.DB
}

func NewGormProfileRepository(db *gorm.DB) *GormProfileRepository {
	return &GormProfileRepository{db: db}
}

func (r *GormProfileRepository) Upsert(profile *domain.Profile) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(profile).Error
}

func (r *GormProfileRepository) GetByUserID(userID uint) (*domain.Profile, error) {
	profile := &domain.Profile{}
	if err := r.db.Where("user_id = ?", userID).First(profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

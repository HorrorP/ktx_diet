package repository

import (
	"time"

	"ktx-diet/internal/domain"

	"gorm.io/gorm"
)

type GormWeeklyPlanRepository struct {
	db *gorm.DB
}

func NewGormWeeklyPlanRepository(db *gorm.DB) *GormWeeklyPlanRepository {
	return &GormWeeklyPlanRepository{db: db}
}

func (r *GormWeeklyPlanRepository) Create(plan *domain.WeeklyPlan) error {
	return r.db.Create(plan).Error
}

func (r *GormWeeklyPlanRepository) GetLatestByUserID(userID uint) (*domain.WeeklyPlan, error) {
	plan := &domain.WeeklyPlan{}
	if err := r.db.Where("user_id = ?", userID).Order("start_date desc").First(plan).Error; err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *GormWeeklyPlanRepository) Update(plan *domain.WeeklyPlan) error {
	return r.db.Save(plan).Error
}

func (r *GormWeeklyPlanRepository) GetByUserIDAndStartDate(userID uint, startDate time.Time) (*domain.WeeklyPlan, error) {
	plan := &domain.WeeklyPlan{}
	if err := r.db.Where("user_id = ? AND start_date = ?", userID, startDate).First(plan).Error; err != nil {
		return nil, err
	}
	return plan, nil
}

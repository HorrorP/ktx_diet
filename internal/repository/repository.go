package repository

import (
	"time"

	"ktx-diet/internal/domain"
)

type UserRepository interface {
	GetOrCreateByTelegramID(telegramID int64) (*domain.User, error)
}

type ProfileRepository interface {
	Upsert(profile *domain.Profile) error
	GetByUserID(userID uint) (*domain.Profile, error)
}

type WeeklyPlanRepository interface {
	Create(plan *domain.WeeklyPlan) error
	GetLatestByUserID(userID uint) (*domain.WeeklyPlan, error)
	Update(plan *domain.WeeklyPlan) error
	GetByUserIDAndStartDate(userID uint, startDate time.Time) (*domain.WeeklyPlan, error)
}

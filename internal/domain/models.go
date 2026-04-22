package domain

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TelegramID int64     `gorm:"uniqueIndex;not null" json:"telegram_id" binding:"required"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Profile    Profile   `json:"profile"`
}

type Profile struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Weight         float64   `gorm:"not null" json:"weight" binding:"required,gt=0"`
	Height         float64   `gorm:"not null" json:"height" binding:"required,gt=0"`
	Age            int       `gorm:"not null" json:"age" binding:"required,gte=10,lte=120"`
	Gender         string    `gorm:"size:16;not null" json:"gender" binding:"required,oneof=male female"`
	ActivityLevel  string    `gorm:"size:24;not null" json:"activity_level" binding:"required,oneof=sedentary light moderate active very_active"`
	TargetCalories int       `gorm:"not null" json:"target_calories"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type WeeklyPlan struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	StartDate time.Time      `gorm:"index;not null" json:"start_date"`
	Content   datatypes.JSON `gorm:"type:jsonb;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

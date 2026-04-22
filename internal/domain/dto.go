package domain

import "time"

type MacroOptions struct {
	WeightLoss  int `json:"weight_loss"`
	Maintenance int `json:"maintenance"`
	Bulking     int `json:"bulking"`
}

type UpsertProfileRequest struct {
	TelegramID    int64   `json:"telegram_id" binding:"required"`
	Weight        float64 `json:"weight" binding:"required,gt=0"`
	Height        float64 `json:"height" binding:"required,gt=0"`
	Age           int     `json:"age" binding:"required,gte=10,lte=120"`
	Gender        string  `json:"gender" binding:"required,oneof=male female"`
	ActivityLevel string  `json:"activity_level" binding:"required,oneof=sedentary light moderate active very_active"`
}

type UpsertProfileResponse struct {
	Profile Profile      `json:"profile"`
	Macros  MacroOptions `json:"macros"`
}

type GeneratePlanRequest struct {
	TelegramID int64  `json:"telegram_id" binding:"required"`
	Goal       string `json:"goal" binding:"required,oneof=weight_loss maintenance bulking"`
	StartDate  string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
}

type GeneratePlanResponse struct {
	Plan WeeklyPlan `json:"plan"`
}

type ReplaceMealRequest struct {
	TelegramID int64  `json:"telegram_id" binding:"required"`
	DayIndex   int    `json:"day_index" binding:"required,gte=0,lte=6"`
	MealType   string `json:"meal_type" binding:"required,oneof=breakfast lunch dinner"`
	Goal       string `json:"goal" binding:"required,oneof=weight_loss maintenance bulking"`
}

type TelegramUser struct {
	ID int64 `json:"id"`
}

type TelegramInitDataPayload struct {
	User TelegramUser `json:"user"`
}

type ServiceContext struct {
	User      User
	Profile   Profile
	StartDate time.Time
}

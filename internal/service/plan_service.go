package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"ktx-diet/internal/domain"
	"ktx-diet/internal/repository"

	"gorm.io/datatypes"
)

type PlanService interface {
	GeneratePlan(ctx context.Context, req domain.GeneratePlanRequest) (*domain.GeneratePlanResponse, error)
	ReplaceMeal(ctx context.Context, req domain.ReplaceMealRequest) (*domain.WeeklyPlan, error)
}

type planService struct {
	userRepo    repository.UserRepository
	profileRepo repository.ProfileRepository
	planRepo    repository.WeeklyPlanRepository
	gemini      GeminiService
}

func NewPlanService(
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	planRepo repository.WeeklyPlanRepository,
	gemini GeminiService,
) PlanService {
	return &planService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		planRepo:    planRepo,
		gemini:      gemini,
	}
}

func (s *planService) GeneratePlan(ctx context.Context, req domain.GeneratePlanRequest) (*domain.GeneratePlanResponse, error) {
	user, err := s.userRepo.GetOrCreateByTelegramID(req.TelegramID)
	if err != nil {
		return nil, err
	}

	profile, err := s.profileRepo.GetByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	targetCalories, err := caloriesByGoal(profile.TargetCalories, req.Goal)
	if err != nil {
		return nil, err
	}

	content, err := s.gemini.GenerateWeeklyPlan(ctx, targetCalories, req.Goal)
	if err != nil {
		return nil, err
	}

	startDate := time.Now().UTC().Truncate(24 * time.Hour)
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		startDate = parsed.UTC()
	}

	plan := &domain.WeeklyPlan{
		UserID:    user.ID,
		StartDate: startDate,
		Content:   datatypes.JSON(content),
	}
	if err := s.planRepo.Create(plan); err != nil {
		return nil, err
	}

	return &domain.GeneratePlanResponse{Plan: *plan}, nil
}

func (s *planService) ReplaceMeal(ctx context.Context, req domain.ReplaceMealRequest) (*domain.WeeklyPlan, error) {
	user, err := s.userRepo.GetOrCreateByTelegramID(req.TelegramID)
	if err != nil {
		return nil, err
	}
	profile, err := s.profileRepo.GetByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	targetCalories, err := caloriesByGoal(profile.TargetCalories, req.Goal)
	if err != nil {
		return nil, err
	}

	plan, err := s.planRepo.GetLatestByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	var parsed map[string]any
	if err := json.Unmarshal(plan.Content, &parsed); err != nil {
		return nil, err
	}

	daysRaw, ok := parsed["days"].([]any)
	if !ok || req.DayIndex >= len(daysRaw) {
		return nil, errors.New("invalid day_index")
	}

	dayObj, ok := daysRaw[req.DayIndex].(map[string]any)
	if !ok {
		return nil, errors.New("invalid day object")
	}

	existingMeal, exists := dayObj[req.MealType]
	if !exists {
		return nil, errors.New("meal_type not found in plan")
	}

	replacementJSON, err := s.gemini.GenerateMealReplacement(ctx, req.Goal, targetCalories, req.MealType, existingMeal)
	if err != nil {
		return nil, err
	}

	var replacement map[string]any
	if err := json.Unmarshal(replacementJSON, &replacement); err != nil {
		return nil, err
	}

	dayObj[req.MealType] = replacement
	daysRaw[req.DayIndex] = dayObj
	parsed["days"] = daysRaw

	updated, err := json.Marshal(parsed)
	if err != nil {
		return nil, err
	}

	plan.Content = datatypes.JSON(updated)
	if err := s.planRepo.Update(plan); err != nil {
		return nil, err
	}

	return plan, nil
}

func caloriesByGoal(maintenance int, goal string) (int, error) {
	switch goal {
	case "weight_loss":
		return int(float64(maintenance) * 0.8), nil
	case "maintenance":
		return maintenance, nil
	case "bulking":
		return int(float64(maintenance) * 1.15), nil
	default:
		return 0, errors.New("invalid goal")
	}
}

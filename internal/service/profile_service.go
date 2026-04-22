package service

import (
	"ktx-diet/internal/domain"
	"ktx-diet/internal/repository"
)

type ProfileService interface {
	UpsertProfile(req domain.UpsertProfileRequest) (*domain.UpsertProfileResponse, error)
}

type profileService struct {
	userRepo      repository.UserRepository
	profileRepo   repository.ProfileRepository
	macroService  MacroService
}

func NewProfileService(
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	macroService MacroService,
) ProfileService {
	return &profileService{
		userRepo:     userRepo,
		profileRepo:  profileRepo,
		macroService: macroService,
	}
}

func (s *profileService) UpsertProfile(req domain.UpsertProfileRequest) (*domain.UpsertProfileResponse, error) {
	user, err := s.userRepo.GetOrCreateByTelegramID(req.TelegramID)
	if err != nil {
		return nil, err
	}

	macros, err := s.macroService.CalculateMacros(req.Weight, req.Height, req.Age, req.Gender, req.ActivityLevel)
	if err != nil {
		return nil, err
	}

	profile := &domain.Profile{
		UserID:         user.ID,
		Weight:         req.Weight,
		Height:         req.Height,
		Age:            req.Age,
		Gender:         req.Gender,
		ActivityLevel:  req.ActivityLevel,
		TargetCalories: macros.Maintenance,
	}

	if err := s.profileRepo.Upsert(profile); err != nil {
		return nil, err
	}

	return &domain.UpsertProfileResponse{
		Profile: *profile,
		Macros:  macros,
	}, nil
}

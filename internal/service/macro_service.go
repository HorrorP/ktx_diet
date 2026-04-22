package service

import (
	"errors"
	"math"

	"ktx-diet/internal/domain"
)

type MacroService interface {
	CalculateMacros(weight, height float64, age int, gender, activityLevel string) (domain.MacroOptions, error)
}

type macroService struct{}

func NewMacroService() MacroService {
	return &macroService{}
}

func (s *macroService) CalculateMacros(weight, height float64, age int, gender, activityLevel string) (domain.MacroOptions, error) {
	genderFactor := 0.0
	switch gender {
	case "male":
		genderFactor = 5
	case "female":
		genderFactor = -161
	default:
		return domain.MacroOptions{}, errors.New("invalid gender")
	}

	activityMultiplier, ok := map[string]float64{
		"sedentary":   1.2,
		"light":       1.375,
		"moderate":    1.55,
		"active":      1.725,
		"very_active": 1.9,
	}[activityLevel]
	if !ok {
		return domain.MacroOptions{}, errors.New("invalid activity level")
	}

	bmr := 10*weight + 6.25*height - 5*float64(age) + genderFactor
	maintenance := int(math.Round(bmr * activityMultiplier))

	return domain.MacroOptions{
		WeightLoss:  int(math.Round(float64(maintenance) * 0.8)),
		Maintenance: maintenance,
		Bulking:     int(math.Round(float64(maintenance) * 1.15)),
	}, nil
}

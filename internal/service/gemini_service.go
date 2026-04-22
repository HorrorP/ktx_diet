package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiService interface {
	GenerateWeeklyPlan(ctx context.Context, targetCalories int, goal string) (json.RawMessage, error)
	GenerateMealReplacement(ctx context.Context, goal string, targetCalories int, mealType string, existingMeal any) (json.RawMessage, error)
}

type geminiService struct {
	model *genai.GenerativeModel
}

func NewGeminiService(ctx context.Context, apiKey string) (GeminiService, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.ResponseMIMEType = "application/json"
	return &geminiService{model: model}, nil
}

func (s *geminiService) GenerateWeeklyPlan(ctx context.Context, targetCalories int, goal string) (json.RawMessage, error) {
	prompt := fmt.Sprintf(`You are a nutrition assistant.
Return STRICT JSON only with this schema:
{
  "days": [
    {
      "day_index": 0,
      "breakfast": {"name":"", "grams": 0, "calories": 0, "recipe": ""},
      "lunch": {"name":"", "grams": 0, "calories": 0, "recipe": ""},
      "dinner": {"name":"", "grams": 0, "calories": 0, "recipe": ""}
    }
  ]
}
Generate exactly 7 days (day_index 0..6), each with breakfast, lunch, dinner.
Each meal must include grams and a short recipe.
Goal: %s.
Daily target calories around %d (+/- 5%%).`, goal, targetCalories)

	out, err := s.generate(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *geminiService) GenerateMealReplacement(ctx context.Context, goal string, targetCalories int, mealType string, existingMeal any) (json.RawMessage, error) {
	existingMealJSON, _ := json.Marshal(existingMeal)
	prompt := fmt.Sprintf(`You are a nutrition assistant.
Return STRICT JSON only with this schema:
{"name":"", "grams": 0, "calories": 0, "recipe": ""}
Generate exactly one %s replacement meal with similar calories/macros to this meal:
%s
Goal: %s.
Target daily calories: %d.
Include grams and a short recipe.`, mealType, string(existingMealJSON), goal, targetCalories)

	out, err := s.generate(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *geminiService) generate(ctx context.Context, prompt string) (json.RawMessage, error) {
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("empty gemini response")
	}

	var sb strings.Builder
	for _, p := range resp.Candidates[0].Content.Parts {
		sb.WriteString(fmt.Sprintf("%v", p))
	}
	content := strings.TrimSpace(sb.String())
	if !json.Valid([]byte(content)) {
		return nil, fmt.Errorf("gemini returned non-json content")
	}
	return json.RawMessage(content), nil
}

package main

import (
	"context"
	"log"
	"os"
	"time"

	"ktx-diet/internal/domain"
	"ktx-diet/internal/handler"
	"ktx-diet/internal/handler/middleware"
	"ktx-diet/internal/repository"
	"ktx-diet/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {


	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	db := mustInitDB()
	mustMigrate(db)

	geminiKey := mustEnv("GEMINI_API_KEY")
	telegramBotToken := mustEnv("TELEGRAM_BOT_TOKEN")
	port := envOrDefault("PORT", "8080")

	userRepo := repository.NewGormUserRepository(db)
	profileRepo := repository.NewGormProfileRepository(db)
	planRepo := repository.NewGormWeeklyPlanRepository(db)

	macroService := service.NewMacroService()
	profileService := service.NewProfileService(userRepo, profileRepo, macroService)

	geminiService, err := service.NewGeminiService(context.Background(), geminiKey)
	if err != nil {
		log.Fatalf("failed to initialize gemini client: %v", err)
	}
	planService := service.NewPlanService(userRepo, profileRepo, planRepo, geminiService)

	profileHandler := handler.NewProfileHandler(profileService)
	planHandler := handler.NewPlanHandler(planService)

	router := gin.Default()
	router.GET("/health", handler.Health)

	api := router.Group("/api/v1")
	api.Use(middleware.TelegramInitDataAuth(telegramBotToken))
	{
		api.POST("/profile", profileHandler.UpsertProfile)
		api.POST("/plan/generate", planHandler.GeneratePlan)
		api.PATCH("/plan/replace", planHandler.ReplaceMeal)
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func mustInitDB() *gorm.DB {
	dsn := mustEnv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql db: %v", err)
	}

	// Set safe defaults for API workloads; tune via env if needed.
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db
}

func mustMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&domain.User{}, &domain.Profile{}, &domain.WeeklyPlan{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env: %s", key)
	}
	return v
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

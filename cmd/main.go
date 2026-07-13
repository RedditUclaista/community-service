package main

import (
	"fmt"
	"github.com/labstack/echo/v5"
	"github.com/RedditUclaista/community-service/internal/config"
	"github.com/RedditUclaista/community-service/internal/database"
	"github.com/RedditUclaista/community-service/internal/usecases"
	deliveryhttp "github.com/RedditUclaista/community-service/internal/delivery/http"
)

func main() {
	cfg := config.LoadConfig()
	
	db, err := database.NewConnection(cfg)
	if err != nil {
		panic("Failed to connect to DB: " + err.Error())
	}
	
	uc := usecases.NewCommunityUseCase(db)
	h := deliveryhttp.NewCommunityHandler(uc)
	
	app := echo.New()
	deliveryhttp.SetupRoutes(app, h, cfg.JWTSecretKey)
	
	fmt.Printf("Starting community-service on port %s...\n", cfg.AppPort)
	if err := app.Start(":" + cfg.AppPort); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

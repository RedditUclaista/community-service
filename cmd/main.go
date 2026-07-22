package main

import (
	"context"
	"fmt"
	"log"

	"github.com/RedditUclaista/community-service/internal/config"
	"github.com/RedditUclaista/community-service/internal/consumer"
	"github.com/RedditUclaista/community-service/internal/database"
	deliveryhttp "github.com/RedditUclaista/community-service/internal/delivery/http"
	"github.com/RedditUclaista/community-service/internal/usecases"
	"github.com/labstack/echo/v5"
)

func main() {
	cfg := config.LoadConfig()

	ctx := context.Background()
	pool, err := database.NewConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pool.Close()

	commRepo := database.NewCommunityRepository(pool)
	memberRepo := database.NewMemberRepository(pool)
	outboxRepo := database.NewOutboxRepository(pool)

	commUc := usecases.NewCommunityUseCase(pool, commRepo, outboxRepo, memberRepo)
	memberUc := usecases.NewMemberUseCase(pool, memberRepo, outboxRepo)

	userRepo := database.NewUserRepository(pool)
	userUc := usecases.NewUserUseCase(userRepo)

	c, err := consumer.NewConsumer(cfg.MQURL, cfg.MQVHost, userUc)
	if err != nil {
		log.Fatalf("Failed to initialize consumer: %v", err)
	}
	if err := c.Start(ctx); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	defer c.Close()

	h := deliveryhttp.NewCommunityHandler(commUc, memberUc)

	app := echo.New()
	deliveryhttp.SetupRoutes(app, h, cfg.JWTSecretKey)

	fmt.Printf("Starting community-service on port %s...\n", cfg.AppPort)
	if err := app.Start(":" + cfg.AppPort); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

package http

import (
	"github.com/RedditUclaista/community-service/internal/middleware"
	"github.com/labstack/echo/v5"
)

func SetupRoutes(app *echo.Echo, h *CommunityHandler, jwtSecret string) {
	v1 := app.Group("/api/v1")

	jwtMd := middleware.JWTMiddleware(jwtSecret)
	adminOrProf := middleware.RequireRole("ADMIN", "PROFESSOR")

	v1.GET("/communities", h.List, jwtMd)
	v1.POST("/communities", h.Create, jwtMd, adminOrProf)
}

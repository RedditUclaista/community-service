package http

import (
	"github.com/RedditUclaista/community-service/internal/middleware"
	"github.com/labstack/echo/v5"
)

func SetupRoutes(app *echo.Echo, h *CommunityHandler, jwtSecret string) {
	api := app.Group("/api/community")

	jwtMd := middleware.JWTMiddleware(jwtSecret)

	api.GET("/communities", h.List)
	api.POST("/communities/bulk", h.GetCommunitiesBulk)
	api.POST("/communities", h.Create, jwtMd)
	api.PUT("/communities/:id", h.Update, jwtMd)

	api.POST("/communities/:id/members", h.Join, jwtMd)
	api.DELETE("/communities/:id/members/:user_id", h.Leave, jwtMd)
	api.PATCH("/communities/:id/members/:user_id/role", h.ChangeRole, jwtMd)
	api.GET("/communities/:id/members", h.GetMembers)
	api.GET("/communities/:id/members/:user_id/role", h.GetMemberRole)

	api.GET("/users/:user_id/communities", h.GetUserCommunities)
}

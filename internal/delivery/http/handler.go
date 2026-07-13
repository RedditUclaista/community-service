package http

import (
	"net/http"
	"github.com/labstack/echo/v5"
	"github.com/RedditUclaista/community-service/internal/usecases"
)

type CommunityHandler struct {
	uc *usecases.CommunityUseCase
}

func NewCommunityHandler(uc *usecases.CommunityUseCase) *CommunityHandler {
	return &CommunityHandler{uc: uc}
}

type CreateReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *CommunityHandler) Create(c *echo.Context) error {
	req := new(CreateReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid format"})
	}
	
	createdBy := c.Get("user_email").(string)
	
	comm, err := h.uc.Create(req.Name, req.Description, createdBy)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, comm)
}

func (h *CommunityHandler) List(c *echo.Context) error {
	list, err := h.uc.ListActive()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

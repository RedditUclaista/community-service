package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RedditUclaista/community-service/internal/dto"
	"github.com/RedditUclaista/community-service/internal/entities"
	"github.com/RedditUclaista/community-service/internal/usecases"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type CommunityHandler struct {
	commUc   *usecases.CommunityUseCase
	memberUc *usecases.MemberUseCase
}

func NewCommunityHandler(commUc *usecases.CommunityUseCase, memberUc *usecases.MemberUseCase) *CommunityHandler {
	return &CommunityHandler{commUc: commUc, memberUc: memberUc}
}

func (h *CommunityHandler) Create(c *echo.Context) error {
	var req dto.CreateCommunityReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid format"})
	}

	fmt.Println(req)

	userIDStr, ok := c.Get("user_id").(string)
	fmt.Println(userIDStr)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id"})
	}

	comm, err := h.commUc.Create(c.Request().Context(), req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, comm)
}

func (h *CommunityHandler) Update(c *echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req dto.UpdateCommunityReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid format"})
	}

	comm, err := h.commUc.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if comm == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, comm)
}

func (h *CommunityHandler) GetCommunitiesBulk(c *echo.Context) error {
	var req dto.RequestBulkCommunities
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	communities, err := h.commUc.GetCommunitiesBulk(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var res dto.CommunityListRes
	for _, comm := range communities {
		res.Communities = append(res.Communities, dto.CommunityRes{
			ID:          comm.ID,
			Name:        comm.Name,
			Description: comm.Description,
			Rules:       comm.Rules,
			BannerURL:   comm.BannerURL,
			ProfileURL:  comm.ProfileURL,
			CreatedBy:   comm.CreatedBy,
		})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *CommunityHandler) List(c *echo.Context) error {
	query := c.QueryParam("query")
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10
	offset := 0
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	communities, err := h.commUc.SearchPaginated(c.Request().Context(), query, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// map to DTO
	res := make([]dto.CommunityRes, len(communities))
	for i, comm := range communities {
		res[i] = dto.CommunityRes{
			ID:          comm.ID,
			Name:        comm.Name,
			Description: comm.Description,
			Rules:       comm.Rules,
			BannerURL:   comm.BannerURL,
			ProfileURL:  comm.ProfileURL,
			CreatedBy:   comm.CreatedBy,
		}
	}

	return c.JSON(http.StatusOK, dto.CommunityListRes{
		Communities: res,
	})
}

func (h *CommunityHandler) Join(c *echo.Context) error {
	idStr := c.Param("id")
	communityID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid community id"})
	}

	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id"})
	}

	if err := h.memberUc.Join(c.Request().Context(), communityID, userID, entities.RoleMember); err != nil {
		// Could check for pgx conflict error here to return 409
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusCreated)
}

func (h *CommunityHandler) Leave(c *echo.Context) error {
	idStr := c.Param("id")
	communityID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid community id"})
	}

	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing user id"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id"})
	}

	if err := h.memberUc.Leave(c.Request().Context(), communityID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h *CommunityHandler) ChangeRole(c *echo.Context) error {
	idStr := c.Param("id")
	communityID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid community id"})
	}

	targetUserIDStr := c.Param("user_id")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	var req dto.ChangeRoleReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid format"})
	}

	if err := h.memberUc.ChangeRole(c.Request().Context(), communityID, targetUserID, req.Role); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h *CommunityHandler) GetUserCommunities(c *echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	communities, err := h.memberUc.GetUserCommunities(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	res := make([]dto.CommunityRes, len(communities))
	for i, comm := range communities {
		res[i] = dto.CommunityRes{
			ID:          comm.ID,
			Name:        comm.Name,
			Description: comm.Description,
			Rules:       comm.Rules,
			BannerURL:   comm.BannerURL,
			ProfileURL:  comm.ProfileURL,
			CreatedBy:   comm.CreatedBy,
			Role:        comm.Role,
		}
	}

	return c.JSON(http.StatusOK, dto.CommunityListRes{
		Communities: res,
	})
}

func (h *CommunityHandler) GetMembers(c *echo.Context) error {
	idStr := c.Param("id")
	communityID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid community id"})
	}

	members, err := h.memberUc.GetMembers(c.Request().Context(), communityID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	res := make([]dto.MemberRes, len(members))
	for i, m := range members {
		res[i] = dto.MemberRes{
			UserID:   m.UserID.String(),
			Role:     string(m.Role),
			JoinedAt: m.JoinedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return c.JSON(http.StatusOK, res)
}

func (h *CommunityHandler) GetMemberRole(c *echo.Context) error {
	idStr := c.Param("id")
	communityID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid community id"})
	}

	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	member, err := h.memberUc.GetMemberRole(c.Request().Context(), communityID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if member == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "member not found"})
	}

	return c.JSON(http.StatusOK, dto.MemberRoleRes{
		Role: string(member.Role),
	})
}

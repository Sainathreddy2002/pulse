package handler

import (
	"errors"
	"net/http"
	"pulse/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	service *service.FollowService
}

func NewFollowHandler(s *service.FollowService) *FollowHandler {
	return &FollowHandler{service: s}
}

func (h *FollowHandler) FollowUser(c *gin.Context) {
	followingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	followerID := c.GetInt64("userID")
	if err := h.service.FollowUser(followingID, followerID); err != nil {
		switch {
		case errors.Is(err, service.ErrCantFollowYourself):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrUserOrFollowerNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrAlreadyFollowing):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "followed"})
}

func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	followingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	followerID := c.GetInt64("userID")
	if err := h.service.UnfollowUser(followingID, followerID); err != nil {
		switch {
		case errors.Is(err, service.ErrUserOrFollowerNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
}

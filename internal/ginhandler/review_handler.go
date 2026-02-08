// internal/ginhandler/reviewhandler.go
package ginhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/AlikhanF2006/Final_project/internal/middleware"
	"github.com/AlikhanF2006/Final_project/internal/postgres/dto"
	"github.com/AlikhanF2006/Final_project/internal/service"
	"github.com/AlikhanF2006/Final_project/model"
)

type ReviewHandler struct {
	reviewSvc *service.ReviewService
}

func NewReviewHandler(reviewSvc *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewSvc: reviewSvc}
}

func (h *ReviewHandler) AddReview(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	var req dto.AddReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	userID := c.GetInt(middleware.UserIDKey)

	created, err := h.reviewSvc.AddReview(movieID, model.Review{UserID: userID, Score: req.Score, Text: req.Text})
	if err != nil {
		if err == service.ErrBadReviewData {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *ReviewHandler) GetReviews(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	revs, err := h.reviewSvc.ListReviews(movieID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, revs)
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	var req dto.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	userID := c.GetInt(middleware.UserIDKey)

	if err := h.reviewSvc.UpdateReview(movieID, userID, req.Score); err != nil {
		if err == service.ErrBadReviewData {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot update review"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	userID := c.GetInt(middleware.UserIDKey)

	if err := h.reviewSvc.DeleteReview(movieID, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete review"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ReviewHandler) AdminDeleteReview(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("review_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	role := c.GetString(middleware.UserRoleKey)
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}

	if err := h.reviewSvc.DeleteReviewByID(reviewID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

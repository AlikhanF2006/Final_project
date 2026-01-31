package ginhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"final_project/internal/service"
	"final_project/model"
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

	var req model.Review
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	created, err := h.reviewSvc.AddReview(movieID, req)
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

package ginhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"final_project/internal/service"
	"final_project/model"
)

type MovieHandler struct {
	movieSvc *service.MovieService
}

func NewMovieHandler(movieSvc *service.MovieService) *MovieHandler {
	return &MovieHandler{movieSvc: movieSvc}
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var req model.Movie
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	created, err := h.movieSvc.CreateMovie(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *MovieHandler) GetMovies(c *gin.Context) {
	c.JSON(http.StatusOK, h.movieSvc.ListMovies())
}

func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	m, err := h.movieSvc.GetMovie(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, m)
}

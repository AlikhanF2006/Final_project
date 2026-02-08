package ginhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/AlikhanF2006/Final_project/internal/service"
	"github.com/AlikhanF2006/Final_project/model"
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

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req model.Movie
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	updated, err := h.movieSvc.UpdateMovie(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.movieSvc.DeleteMovie(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *MovieHandler) GetPopularFromTMDB(c *gin.Context) {
	movies, err := h.movieSvc.GetPopularFromTMDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, movies)
}

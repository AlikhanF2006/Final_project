package service

import (
	"errors"
	"strings"

	"github.com/AlikhanF2006/Final_project/internal/postgres"
	"github.com/AlikhanF2006/Final_project/model"
)

var (
	ErrBadMovieData = errors.New("invalid movie data")
)

type MovieService struct {
	movieRepo *postgres.MovieRepository
}

func NewMovieService(movieRepo *postgres.MovieRepository) *MovieService {
	return &MovieService{movieRepo: movieRepo}
}

func (s *MovieService) CreateMovie(m model.Movie) (model.Movie, error) {
	m.Title = strings.TrimSpace(m.Title)
	if m.Title == "" || m.Year <= 0 {
		return model.Movie{}, ErrBadMovieData
	}

	m.ID = 0
	m.Rating = 0

	return s.movieRepo.Create(m), nil
}

func (s *MovieService) ListMovies() []model.Movie {
	return s.movieRepo.GetAll()
}

func (s *MovieService) GetMovie(id int) (model.Movie, error) {
	return s.movieRepo.GetByID(id)
}

func (s *MovieService) SetRating(movieID int, rating float64) error {
	return s.movieRepo.SetRating(movieID, rating)
}

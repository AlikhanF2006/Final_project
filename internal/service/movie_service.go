package service

import (
	"errors"
	"strings"

	"github.com/AlikhanF2006/Final_project/internal/postgres"
	"github.com/AlikhanF2006/Final_project/model"
)

var ErrBadMovieData = errors.New("invalid movie data")

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
	return s.movieRepo.Create(m), nil
}

func (s *MovieService) ListMovies() []model.Movie {
	return s.movieRepo.GetAll()
}

func (s *MovieService) GetMovie(id int) (model.Movie, error) {
	return s.movieRepo.GetByID(id)
}

func (s *MovieService) UpdateMovie(id int, upd model.Movie) (model.Movie, error) {
	existing, err := s.movieRepo.GetByID(id)
	if err != nil {
		return model.Movie{}, err
	}

	if strings.TrimSpace(upd.Title) != "" {
		existing.Title = strings.TrimSpace(upd.Title)
	}
	if upd.Year > 0 {
		existing.Year = upd.Year
	}
	if upd.Description != "" {
		existing.Description = upd.Description
	}

	return s.movieRepo.Update(existing)
}

func (s *MovieService) DeleteMovie(id int) error {
	return s.movieRepo.Delete(id)
}

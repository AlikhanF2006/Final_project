package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/AlikhanF2006/Final_project/internal/postgres"
	"github.com/AlikhanF2006/Final_project/internal/tmdb"
	"github.com/AlikhanF2006/Final_project/model"
)

var ErrBadMovieData = errors.New("invalid movie data")

type MovieService struct {
	movieRepo  *postgres.MovieRepository
	tmdbClient *tmdb.Client
}

func NewMovieService(
	movieRepo *postgres.MovieRepository,
	tmdbClient *tmdb.Client,
) *MovieService {
	return &MovieService{
		movieRepo:  movieRepo,
		tmdbClient: tmdbClient,
	}
}

func (s *MovieService) CreateMovie(m model.Movie) (model.Movie, error) {
	m.Title = strings.TrimSpace(m.Title)
	if m.Title == "" || m.Year <= 0 {
		return model.Movie{}, ErrBadMovieData
	}
	return s.movieRepo.Create(m)
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

func (s *MovieService) Search(title string, year int) []model.Movie {
	all := s.movieRepo.GetAll()
	result := make([]model.Movie, 0)

	title = strings.ToLower(strings.TrimSpace(title))

	for _, m := range all {
		if title != "" && !strings.Contains(strings.ToLower(m.Title), title) {
			continue
		}
		if year > 0 && m.Year != year {
			continue
		}
		result = append(result, m)
	}

	return result
}

func (s *MovieService) GetPopularFromTMDB() ([]model.Movie, error) {
	moviesDTO, err := s.tmdbClient.GetPopularMovies()
	if err != nil {
		return nil, err
	}

	var result []model.Movie

	for _, m := range moviesDTO {
		year := 0
		if len(m.ReleaseDate) >= 4 {
			if y, err := strconv.Atoi(m.ReleaseDate[:4]); err == nil {
				year = y
			}
		}

		if s.movieRepo.ExistsByTMDBID(m.ID) {
			existing, err := s.movieRepo.GetByTMDBID(m.ID)
			if err == nil {
				result = append(result, existing)
			}
			continue
		}

		created, err := s.movieRepo.Create(model.Movie{
			TMDBID:      m.ID,
			Title:       m.Title,
			Description: m.Overview,
			Year:        year,
			Rating:      0,
		})
		if err == nil {
			result = append(result, created)
		}
	}

	return result, nil
}

func (s *MovieService) GetMovieWithTrailer(tmdbID int) (map[string]any, error) {
	movie, err := s.tmdbClient.GetMovie(tmdbID)
	if err != nil {
		return nil, err
	}

	trailerKey, _ := s.tmdbClient.GetTrailerKey(tmdbID)

	result := map[string]any{
		"id":           movie.ID,
		"title":        movie.Title,
		"description":  movie.Overview,
		"release_date": movie.ReleaseDate,
		"trailer_url":  "",
	}

	if trailerKey != "" {
		result["trailer_url"] = "https://www.youtube.com/watch?v=" + trailerKey
	}

	return result, nil
}

func (s *MovieService) SearchMovies(title string, year int) ([]model.Movie, error) {
	return s.movieRepo.Search(title, year)
}

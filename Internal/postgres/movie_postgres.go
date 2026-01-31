package postgres

import (
	"errors"
	"sync"

	"github.com/AlikhanF2006/Final_project/model"
)

var (
	ErrMovieNotFound = errors.New("movie not found")
)

type MovieRepository struct {
	mu     sync.RWMutex
	nextID int
	movies map[int]model.Movie
}

func NewMovieRepository() *MovieRepository {
	return &MovieRepository{
		nextID: 1,
		movies: make(map[int]model.Movie),
	}
}

func (r *MovieRepository) Create(m model.Movie) model.Movie {
	r.mu.Lock()
	defer r.mu.Unlock()

	m.ID = r.nextID
	r.nextID++
	r.movies[m.ID] = m
	return m
}

func (r *MovieRepository) GetAll() []model.Movie {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]model.Movie, 0, len(r.movies))
	for _, m := range r.movies {
		out = append(out, m)
	}
	return out
}

func (r *MovieRepository) GetByID(id int) (model.Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.movies[id]
	if !ok {
		return model.Movie{}, ErrMovieNotFound
	}
	return m, nil
}

func (r *MovieRepository) SetRating(movieID int, rating float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, ok := r.movies[movieID]
	if !ok {
		return ErrMovieNotFound
	}
	m.Rating = rating
	r.movies[movieID] = m
	return nil
}

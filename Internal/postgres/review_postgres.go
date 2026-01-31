package postgres

import (
	"sync"
	"time"

	"final_project/model"
)

type ReviewRepository struct {
	mu      sync.RWMutex
	nextID  int
	reviews []model.Review
}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{
		nextID:  1,
		reviews: make([]model.Review, 0),
	}
}

func (r *ReviewRepository) Add(movieID int, review model.Review) model.Review {
	r.mu.Lock()
	defer r.mu.Unlock()

	review.ID = r.nextID
	r.nextID++
	review.MovieID = movieID
	review.CreatedAt = time.Now()

	r.reviews = append(r.reviews, review)
	return review
}

func (r *ReviewRepository) ListByMovieID(movieID int) []model.Review {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]model.Review, 0)
	for _, rev := range r.reviews {
		if rev.MovieID == movieID {
			out = append(out, rev)
		}
	}
	return out
}

package service

import (
	"errors"

	"github.com/AlikhanF2006/Final_project/internal/postgres"
	"github.com/AlikhanF2006/Final_project/model"
)

var (
	ErrBadReviewData = errors.New("invalid review data")
)

type ReviewService struct {
	reviewRepo *postgres.ReviewRepository
	movieRepo  *postgres.MovieRepository
}

func NewReviewService(
	reviewRepo *postgres.ReviewRepository,
	movieRepo *postgres.MovieRepository,
) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		movieRepo:  movieRepo,
	}
}

func (s *ReviewService) AddReview(movieID int, r model.Review) (model.Review, error) {
	if _, err := s.movieRepo.GetByID(movieID); err != nil {
		return model.Review{}, err
	}

	if r.UserID <= 0 || r.Score < 1 || r.Score > 5 {
		return model.Review{}, ErrBadReviewData
	}

	created := s.reviewRepo.Add(movieID, r)

	s.recalculateRating(movieID)

	return created, nil
}

func (s *ReviewService) ListReviews(movieID int) ([]model.Review, error) {
	if _, err := s.movieRepo.GetByID(movieID); err != nil {
		return nil, err
	}
	return s.reviewRepo.ListByMovieID(movieID), nil
}

func (s *ReviewService) recalculateRating(movieID int) {
	revs := s.reviewRepo.ListByMovieID(movieID)
	if len(revs) == 0 {
		_ = s.movieRepo.SetRating(movieID, 0)
		return
	}

	sum := 0
	for _, r := range revs {
		sum += r.Score
	}
	avg := float64(sum) / float64(len(revs))
	_ = s.movieRepo.SetRating(movieID, avg)
}

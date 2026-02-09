package service

import (
	"errors"
	"fmt"

	"github.com/AlikhanF2006/Final_project/internal/postgres"
	"github.com/AlikhanF2006/Final_project/model"
)

var (
	ErrBadReviewData = errors.New("invalid review data")
	ErrForbidden     = errors.New("forbidden")
)

type ReviewService struct {
	reviewRepo *postgres.ReviewRepository
	movieRepo  *postgres.MovieRepository
	ratingCh   chan int
}

func NewReviewService(
	reviewRepo *postgres.ReviewRepository,
	movieRepo *postgres.MovieRepository,
) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		movieRepo:  movieRepo,
		ratingCh:   make(chan int, 10),
	}
}

func (s *ReviewService) StartRatingWorker() {
	go func() {
		for movieID := range s.ratingCh {
			s.recalculateRating(movieID)
		}
	}()
}

func (s *ReviewService) AddReview(movieID int, r model.Review) (model.Review, error) {
	if _, err := s.movieRepo.GetByID(movieID); err != nil {
		return model.Review{}, err
	}

	if r.UserID <= 0 || r.Score < 1 || r.Score > 5 {
		return model.Review{}, ErrBadReviewData
	}

	created, err := s.reviewRepo.Add(movieID, r)
	if err != nil {
		return model.Review{}, err
	}

	s.ratingCh <- movieID
	return created, nil
}

func (s *ReviewService) ListReviews(movieID int) ([]model.Review, error) {
	if _, err := s.movieRepo.GetByID(movieID); err != nil {
		return nil, err
	}
	return s.reviewRepo.ListByMovieID(movieID)
}

func (s *ReviewService) UpdateReview(
	movieID int,
	userID int,
	newScore int,
) error {
	if newScore < 1 || newScore > 5 {
		return ErrBadReviewData
	}

	if err := s.reviewRepo.UpdateByMovieAndUser(
		movieID,
		userID,
		newScore,
	); err != nil {
		return ErrForbidden
	}

	s.ratingCh <- movieID
	return nil
}

func (s *ReviewService) DeleteReview(
	movieID int,
	userID int,
) error {
	if err := s.reviewRepo.DeleteByMovieAndUser(
		movieID,
		userID,
	); err != nil {
		return ErrForbidden
	}

	s.ratingCh <- movieID
	return nil
}

func (s *ReviewService) DeleteReviewByID(reviewID int) error {
	rev, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return fmt.Errorf("review not found")
	}
	if err := s.reviewRepo.DeleteByID(reviewID); err != nil {
		return fmt.Errorf("cannot delete review")
	}
	s.ratingCh <- rev.MovieID
	return nil
}

func (s *ReviewService) recalculateRating(movieID int) {
	revs, err := s.reviewRepo.ListByMovieID(movieID)
	if err != nil || len(revs) == 0 {
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

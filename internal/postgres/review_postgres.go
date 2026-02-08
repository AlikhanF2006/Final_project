package postgres

import (
	"context"

	"github.com/AlikhanF2006/Final_project/model"
	"github.com/AlikhanF2006/Final_project/pkg/db"
)

type ReviewRepository struct{}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{}
}

func (r *ReviewRepository) Add(movieID int, rev model.Review) (model.Review, error) {
	query := `
		INSERT INTO reviews (movie_id, user_id, score)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := db.DB.QueryRow(
		context.Background(),
		query,
		movieID,
		rev.UserID,
		rev.Score,
	).Scan(&rev.ID, &rev.CreatedAt)

	rev.MovieID = movieID
	return rev, err
}

func (r *ReviewRepository) ListByMovieID(movieID int) ([]model.Review, error) {
	query := `
		SELECT id, movie_id, user_id, score, created_at
		FROM reviews
		WHERE movie_id = $1
	`

	rows, err := db.DB.Query(context.Background(), query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	revs := make([]model.Review, 0)
	for rows.Next() {
		var r model.Review
		if err := rows.Scan(
			&r.ID,
			&r.MovieID,
			&r.UserID,
			&r.Score,
			&r.CreatedAt,
		); err != nil {
			return nil, err
		}
		revs = append(revs, r)
	}

	return revs, nil
}

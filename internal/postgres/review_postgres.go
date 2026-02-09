package postgres

import (
	"context"
	"errors"

	"github.com/AlikhanF2006/Final_project/model"
	"github.com/AlikhanF2006/Final_project/pkg/db"
)

type ReviewRepository struct{}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{}
}

func (r *ReviewRepository) Add(movieID int, rev model.Review) (model.Review, error) {
	query := `
		INSERT INTO reviews (movie_id, user_id, score, text)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := db.DB.QueryRow(
		context.Background(),
		query,
		movieID,
		rev.UserID,
		rev.Score,
		rev.Text,
	).Scan(&rev.ID, &rev.CreatedAt)

	rev.MovieID = movieID
	return rev, err
}

func (r *ReviewRepository) ListByMovieID(movieID int) ([]model.Review, error) {
	query := `
		SELECT id, movie_id, user_id, score, text, created_at
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
		var rr model.Review
		if err := rows.Scan(
			&rr.ID,
			&rr.MovieID,
			&rr.UserID,
			&rr.Score,
			&rr.Text,
			&rr.CreatedAt,
		); err != nil {
			return nil, err
		}
		revs = append(revs, rr)
	}

	return revs, nil
}

func (r *ReviewRepository) UpdateByMovieAndUser(
	movieID int,
	userID int,
	score int,
) error {
	cmd, err := db.DB.Exec(
		context.Background(),
		`UPDATE reviews SET score=$1 WHERE movie_id=$2 AND user_id=$3`,
		score,
		movieID,
		userID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *ReviewRepository) DeleteByMovieAndUser(
	movieID int,
	userID int,
) error {
	cmd, err := db.DB.Exec(
		context.Background(),
		`DELETE FROM reviews WHERE movie_id=$1 AND user_id=$2`,
		movieID,
		userID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *ReviewRepository) GetByID(id int) (model.Review, error) {
	var rev model.Review
	query := `SELECT id, movie_id, user_id, score, text, created_at FROM reviews WHERE id=$1`
	err := db.DB.QueryRow(context.Background(), query, id).Scan(
		&rev.ID,
		&rev.MovieID,
		&rev.UserID,
		&rev.Score,
		&rev.Text,
		&rev.CreatedAt,
	)
	if err != nil {
		return model.Review{}, errors.New("not found")
	}
	return rev, nil
}

func (r *ReviewRepository) DeleteByID(id int) error {
	cmd, err := db.DB.Exec(context.Background(), `DELETE FROM reviews WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

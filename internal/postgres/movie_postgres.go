package postgres

import (
	"context"
	"errors"

	"github.com/AlikhanF2006/Final_project/model"
	"github.com/AlikhanF2006/Final_project/pkg/db"
)

var (
	ErrMovieNotFound = errors.New("movie not found")
)

type MovieRepository struct{}

func NewMovieRepository() *MovieRepository {
	return &MovieRepository{}
}

func (r *MovieRepository) Create(m model.Movie) model.Movie {
	query := `
		INSERT INTO movies (title, year, description, rating)
		VALUES ($1, $2, $3, 0)
		RETURNING id, rating
	`

	_ = db.DB.QueryRow(
		context.Background(),
		query,
		m.Title,
		m.Year,
		m.Description,
	).Scan(&m.ID, &m.Rating)

	return m
}

func (r *MovieRepository) GetAll() []model.Movie {
	rows, err := db.DB.Query(
		context.Background(),
		`SELECT id, title, year, description, rating FROM movies`,
	)
	if err != nil {
		return []model.Movie{}
	}
	defer rows.Close()

	var movies []model.Movie

	for rows.Next() {
		var m model.Movie
		if err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Year,
			&m.Description,
			&m.Rating,
		); err == nil {
			movies = append(movies, m)
		}
	}

	return movies
}

func (r *MovieRepository) GetByID(id int) (model.Movie, error) {
	var m model.Movie

	err := db.DB.QueryRow(
		context.Background(),
		`SELECT id, title, year, description, rating FROM movies WHERE id=$1`,
		id,
	).Scan(
		&m.ID,
		&m.Title,
		&m.Year,
		&m.Description,
		&m.Rating,
	)

	if err != nil {
		return model.Movie{}, ErrMovieNotFound
	}

	return m, nil
}

func (r *MovieRepository) SetRating(movieID int, rating float64) error {
	cmd, err := db.DB.Exec(
		context.Background(),
		`UPDATE movies SET rating=$1 WHERE id=$2`,
		rating,
		movieID,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrMovieNotFound
	}

	return nil
}

func (r *MovieRepository) Update(m model.Movie) (model.Movie, error) {
	cmd, err := db.DB.Exec(
		context.Background(),
		`UPDATE movies SET title=$1, year=$2, description=$3 WHERE id=$4`,
		m.Title,
		m.Year,
		m.Description,
		m.ID,
	)

	if err != nil {
		return model.Movie{}, err
	}

	if cmd.RowsAffected() == 0 {
		return model.Movie{}, ErrMovieNotFound
	}

	return m, nil
}

func (r *MovieRepository) Delete(id int) error {
	cmd, err := db.DB.Exec(
		context.Background(),
		`DELETE FROM movies WHERE id=$1`,
		id,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrMovieNotFound
	}

	return nil
}

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

func (r *MovieRepository) SetRating(movieID int, rating float64) error {
	_, err := db.DB.Exec(
		context.Background(),
		`UPDATE movies SET rating = $1 WHERE id = $2`,
		rating,
		movieID,
	)
	return err
}

func (r *MovieRepository) Create(m model.Movie) (model.Movie, error) {
	query := `
		INSERT INTO movies (tmdb_id, title, year, description, rating)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := db.DB.QueryRow(
		context.Background(),
		query,
		m.TMDBID,
		m.Title,
		m.Year,
		m.Description,
		m.Rating,
	).Scan(&m.ID)

	return m, err
}

func (r *MovieRepository) GetAll() []model.Movie {
	rows, err := db.DB.Query(
		context.Background(),
		`SELECT id, tmdb_id, title, year, description, rating FROM movies`,
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
			&m.TMDBID,
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
		`SELECT id, tmdb_id, title, year, description, rating FROM movies WHERE id=$1`,
		id,
	).Scan(
		&m.ID,
		&m.TMDBID,
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

func (r *MovieRepository) GetByTMDBID(tmdbID int) (model.Movie, error) {
	var m model.Movie

	err := db.DB.QueryRow(
		context.Background(),
		`SELECT id, tmdb_id, title, year, description, rating FROM movies WHERE tmdb_id=$1`,
		tmdbID,
	).Scan(
		&m.ID,
		&m.TMDBID,
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

func (r *MovieRepository) ExistsByTMDBID(tmdbID int) bool {
	var id int
	err := db.DB.QueryRow(
		context.Background(),
		`SELECT id FROM movies WHERE tmdb_id=$1`,
		tmdbID,
	).Scan(&id)

	return err == nil
}

func (r *MovieRepository) Update(m model.Movie) (model.Movie, error) {
	cmd, err := db.DB.Exec(
		context.Background(),
		`UPDATE movies SET title=$1, year=$2, description=$3, rating=$4 WHERE id=$5`,
		m.Title,
		m.Year,
		m.Description,
		m.Rating,
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

func (r *MovieRepository) Search(title string, year int) ([]model.Movie, error) {
	query := `
		SELECT id, tmdb_id, title, year, description, rating
		FROM movies
		WHERE ($1 = '' OR LOWER(title) LIKE '%' || LOWER($1) || '%')
		  AND ($2 = 0 OR year = $2)
	`

	rows, err := db.DB.Query(
		context.Background(),
		query,
		title,
		year,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var m model.Movie
		if err := rows.Scan(
			&m.ID,
			&m.TMDBID,
			&m.Title,
			&m.Year,
			&m.Description,
			&m.Rating,
		); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

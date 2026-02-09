package postgres

import (
	"context"
	"errors"

	"github.com/AlikhanF2006/Final_project/model"
	"github.com/AlikhanF2006/Final_project/pkg/db"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(u model.User) (model.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := db.DB.QueryRow(
		context.Background(),
		query,
		u.Username,
		u.Email,
		u.PasswordHash,
		u.Role,
	).Scan(&u.ID, &u.CreatedAt)

	return u, err
}

func (r *UserRepository) GetByEmail(email string) (model.User, error) {
	var u model.User
	query := `
		SELECT id, username, email, password_hash, role, created_at
		FROM users WHERE email=$1
	`
	err := db.DB.QueryRow(context.Background(), query, email).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		return model.User{}, ErrUserNotFound
	}
	return u, nil
}

func (r *UserRepository) GetByID(id int) (model.User, error) {
	var u model.User
	query := `
		SELECT id, username, email, password_hash, role, created_at
		FROM users WHERE id=$1
	`
	err := db.DB.QueryRow(context.Background(), query, id).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		return model.User{}, ErrUserNotFound
	}
	return u, nil
}

func (r *UserRepository) Update(u model.User) (model.User, error) {
	_, err := db.DB.Exec(
		context.Background(),
		`UPDATE users SET username=$1, email=$2 WHERE id=$3`,
		u.Username,
		u.Email,
		u.ID,
	)
	return u, err
}

func (r *UserRepository) UpdatePassword(id int, hash string) error {
	_, err := db.DB.Exec(context.Background(), `UPDATE users SET password_hash=$1 WHERE id=$2`, hash, id)
	return err
}

func (r *UserRepository) Delete(id int) error {
	_, err := db.DB.Exec(
		context.Background(),
		`DELETE FROM users WHERE id=$1`,
		id,
	)
	return err
}

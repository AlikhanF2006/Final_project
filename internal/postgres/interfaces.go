package postgres

import "github.com/AlikhanF2006/Final_project/model"

type MovieRepo interface {
	Create(model.Movie) (model.Movie, error)
	GetAll() []model.Movie
	GetByID(int) (model.Movie, error)
	GetByTMDBID(int) (model.Movie, error)
	ExistsByTMDBID(int) bool
	Update(model.Movie) (model.Movie, error)
	Delete(int) error
	SetRating(int, float64) error
}

type ReviewRepo interface {
	Add(int, model.Review) (model.Review, error)
	ListByMovieID(int) ([]model.Review, error)
	UpdateByMovieAndUser(int, int, int) error
	DeleteByMovieAndUser(int, int) error
	GetByID(int) (model.Review, error)
	DeleteByID(int) error
}

type UserRepo interface {
	Create(model.User) (model.User, error)
	GetByEmail(string) (model.User, error)
	GetByID(int) (model.User, error)
	Update(model.User) (model.User, error)
	UpdatePassword(int, string) error
	Delete(int) error
}

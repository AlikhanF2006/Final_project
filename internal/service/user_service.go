// internal/service/userservice.go
package service

import (
	"errors"
	"time"

	"github.com/AlikhanF2006/Final_project/internal/auth"
	"github.com/AlikhanF2006/Final_project/internal/postgres"
	"github.com/AlikhanF2006/Final_project/internal/postgres/dto"
	"github.com/AlikhanF2006/Final_project/model"
	"golang.org/x/crypto/bcrypt"
)

var ErrBadCredentials = errors.New("invalid credentials")

type UserService struct {
	repo *postgres.UserRepository
}

func NewUserService(r *postgres.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Register(req dto.RegisterDTO) (dto.UserDTO, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         "user",
	}

	created, err := s.repo.Create(user)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return toUserDTO(created), nil
}

func (s *UserService) Login(req dto.LoginDTO) (string, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return "", ErrBadCredentials
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	) != nil {
		return "", ErrBadCredentials
	}

	return auth.GenerateToken(user.ID)
}

func (s *UserService) GetProfile(id int) (dto.UserDTO, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return dto.UserDTO{}, err
	}
	return toUserDTO(u), nil
}

func (s *UserService) UpdateProfile(id int, req dto.UpdateProfileDTO) (dto.UserDTO, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return dto.UserDTO{}, err
	}

	if req.Username != "" {
		u.Username = req.Username
	}
	if req.Email != "" {
		u.Email = req.Email
	}

	updated, err := s.repo.Update(u)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return toUserDTO(updated), nil
}

func (s *UserService) ChangePassword(id int, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("password too short")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	return s.repo.UpdatePassword(id, string(hash))
}

func (s *UserService) DeleteAccount(id int) error {
	return s.repo.Delete(id)
}

func (s *UserService) AdminDeleteUser(id int) error {
	return s.repo.Delete(id)
}

func toUserDTO(u model.User) dto.UserDTO {
	return dto.UserDTO{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

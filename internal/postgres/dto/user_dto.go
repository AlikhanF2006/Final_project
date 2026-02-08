// internal/postgres/dto/user_dto.go
package dto

type UserDTO struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type RegisterDTO struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ChangePasswordDTO struct {
	Password string `json:"password" binding:"required,min=6"`
}

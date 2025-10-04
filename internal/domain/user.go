package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8" hash:"bcrypt"`
	Username  string    `json:"username" validate:"required,min=3,max=20"`
	FirstName string    `json:"first_name" validate:"required,min=3,max=20"`
	LastName  string    `json:"last_name" validate:"required,min=3,max=20"`
	Role      string    `json:"role" default:"user"`
	IsActive  bool      `json:"is_active" default:"true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=20"`
	LastName  string `json:"last_name" validate:"required,min=3,max=20"`
}

type UpdatePasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
}



func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

package types

import (
	"time"
)

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	FullName  string    `gorm:"type:varchar(100);not null" json:"full_name" validate:"required,min=3,max=100"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string    `gorm:"not null" json:"-" validate:"required,min=6,max=100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterPayload struct {
	FullName string `json:"full_name" binding:"required,min=3,max=100" validate:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100" validate:"required,min=6,max=100"`
}

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100" validate:"required,min=6,max=100"`
}

type TokenResponse struct {
	Message               string `json:"message"`
	Token                 string `json:"token"`
	ExpiredAt             string `json:"expired_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiredAt string `json:"refresh_token_expired_at"`
	Role                  string `json:"role"`
}

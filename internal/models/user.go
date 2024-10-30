package models

import "time"

type UserRegistration struct {
	Name     string `json:"name" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type UserProfile struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Description  string    `json:"description"`
	JoinedAt     time.Time `json:"joinedAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	PasswordHash string    `json:"-"`
}

type UserProfileUpdate struct {
	NewName string `json:"name"`
	Email   string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

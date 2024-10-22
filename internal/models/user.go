package models

import "time"

type UserRegistration struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Description string    `json:"description"`
	JoinedAt    time.Time `json:"joinedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

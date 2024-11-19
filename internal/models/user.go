package models

import "time"

type UserRegistration struct {
	Name     string `json:"name" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type UserProfile struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Description    string    `json:"description"`
	JoinedAt       time.Time `json:"joinedAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	AvatarImageURL string    `json:"avatarImageUrl"`
}

type UserProfileUpdate struct {
	NewName string `json:"name" validate:"required,min=3,max=30"`
	Email   string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	NewPassword string `json:"newPassword" validate:"required,min=8,max=50"`
	OldPassword string `json:"oldPassword"`
}

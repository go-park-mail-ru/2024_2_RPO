package models

import "time"

type UserProfile struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Description    string    `json:"description"`
	JoinedAt       time.Time `json:"joinedAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	AvatarImageURL string    `json:"avatarImageUrl"`
}

package models

//go:generate easyjson -all user.go

import "time"

type UserProfile struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	JoinedAt       time.Time      `json:"joinedAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	AvatarImageURL string         `json:"avatarImageUrl"`
	CsatPollDT     time.Time      `json:"-"`
	PollQuestions  []PollQuestion `json:"pollQuestions,omitempty"`
}

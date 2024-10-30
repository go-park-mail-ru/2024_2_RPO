package models

type UserWithPermissions struct {
	User User `json:"user"`
	AddedAt
}

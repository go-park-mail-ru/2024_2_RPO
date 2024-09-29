package models

type CreateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Background  string `json:"background,omitempty"`
}

type Board struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Background  string `json:"background,omitempty"`
	OwnerID     int    `json:"owner_id"`
}

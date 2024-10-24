package models

type BadResponse struct {
	Status int    `json:"status"`
	Text   string `json:"text"`
}

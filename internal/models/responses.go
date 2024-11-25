package models

type BadResponse struct {
	Status int    `json:"status"`
	Text   string `json:"text"`
}

// Если карточка находится на той доске, на которой пользователь есть
type SharedCardFoundResponse struct {
	BoardID int `json:"boardId"`
	CardID  int `json:"cardId"`
}

// Если пользователь не имеет доступа к доске, на которой эта карточка есть
type SharedCardDummyResponse struct {
	BoardName          string       `json:"boardName"`
	BackgroundImageURL string       `json:"backgroundImageUrl"`
	Card               *CardDetails `json:"card"`
}

package models

type BadResponse struct {
	Status int    `json:"status"`
	Text   string `json:"text"`
}

// Если карточка находится на той доске, на которой пользователь есть
type SharedCardFoundResponse struct {
	BoardID int64 `json:"boardId"`
	CardID  int64 `json:"cardId"`
}

// Если пользователь не имеет доступа к доске, на которой эта карточка есть
type SharedCardDummyResponse struct {
	Board *Board       `json:"board"`
	Card  *CardDetails `json:"card"`
}

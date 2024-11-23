package models

type PollResult struct {
	QuestionID int64  `json:"questionId"`
	Type       string `json:"type"`
	Text       string `json:"text"`
	Rating     int64  `json:"rating"`
}

type PollQuestion struct {
	QuestionID   int64  `json:"questionId"`
	QuestionText string `json:"questionText"`
	Type         string `json:"type"`
}

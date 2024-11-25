package models

type RatingResults struct {
	Question string `json:"question"`
	Rating   string `json:"rating"`
}

type AnswerResults struct {
	Question string   `json:"question" `
	Text     []string `json:"text" `
}

type PollQuestion struct {
	QuestionID   int64  `json:"questionId" `
	QuestionText string `json:"questionText" `
	QuestionType string `json:"questionType" `
}

type PollSubmit struct {
	QuestionID   int64   `json:"questionId"`
	QuestionType string  `json:"questionType" `
	Rating       *int    `json:"rating" `
	Text         *string `json:"text" `
}

type PollResults struct {
	RatingResults []RatingResults `json:"ratingResults"`
	TextResults   []AnswerResults `json:"textResults"`
}

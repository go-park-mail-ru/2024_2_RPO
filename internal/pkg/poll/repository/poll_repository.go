package repository

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/pgxiface"
	"context"
	"fmt"
	"time"
)

type PollRepository struct {
	db pgxiface.PgxIface
}

func CreatePollRepository(db pgxiface.PgxIface) *PollRepository {
	return &PollRepository{
		db: db,
	}
}

func (r *PollRepository) SubmitPoll(ctx context.Context, userID int64, PollSubmit *models.PollSubmit) error {
	funcName := "SubmitPoll"
	query := `
	INSERT INTO csat_results (question_id, rating, comment, u_id, created_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP);
	`

	_, err := r.db.Exec(ctx, query, PollSubmit.QuestionID, PollSubmit.Rating, PollSubmit.Text, userID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("SubmitPoll (query): %w", err)
	}

	return nil
}

func (r *PollRepository) GetRatingResults(ctx context.Context) (results []models.RatingResults, err error) {
	funcName := "GetRatingResults"
	query := `
	SELECT cq.question_text, AVG(cr.rating) AS rating FROM csat_results AS cr
	JOIN csat_question AS cq ON cr.question_id = cq.question_id
	WHERE cq.type='answer_rating'
	AND cr.rating IS NOT NULL
	GROUP BY cq.question_id;
	`

	results = make([]models.RatingResults, 0)

	rows, err := r.db.Query(ctx, query)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("GetRatingResults (query): %w", err)
	}

	for rows.Next() {
		result := models.RatingResults{}
		if err := rows.Scan(&result.Question, &result.Rating); err != nil {
			return nil, fmt.Errorf("GetRatingResults (scan): %w", err)
		}
		fmt.Printf("RATING RES: %#v\n", result)
		results = append(results, result)
	}

	return results, nil
}

func (r *PollRepository) GetTextResults(ctx context.Context) (results []models.AnswerResults, err error) {
	funcName := "GetTextResults"
	query := `
	SELECT cr.comment, cq.question_text FROM csat_results AS cr
	JOIN csat_question AS cq ON cr.question_id = cq.question_id
	WHERE cq.type='answer_text'
	AND cr.comment IS NOT NULL
	ORDER BY cq.question_id;
	`

	results = make([]models.AnswerResults, 0)
	rows, err := r.db.Query(ctx, query)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("GetTextResults (query): %w", err)
	}

	for rows.Next() {
		var q, a string
		if err := rows.Scan(&a, &q); err != nil {
			return nil, fmt.Errorf("GetTextResults (scan): %w", err)
		}
		fmt.Println("TEXT RESULT: ", q, "  ", a)
		if len(results) == 0 {
			results = append(results, models.AnswerResults{
				Question: q,
				Text:     []string{a},
			})
		} else {
			if results[len(results)-1].Question == q {
				results[len(results)-1].Text = append(results[len(results)-1].Text, a)
			} else {
				results = append(results, models.AnswerResults{
					Question: q,
					Text:     []string{a},
				})
			}
		}
	}

	return results, nil
}

func (r *PollRepository) SetNextPollDT(ctx context.Context, userID int64) error {
	funcName := "SetNextPollDate"
	query := `
	UPDATE "user" SET csat_poll_dt=(CURRENT_TIMESTAMP+$2) WHERE u_id=$1;
	`

	_, err := r.db.Exec(ctx, query, userID, 24*7*time.Hour)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("SetNextPollDate (query): %w", err)
	}

	return nil
}

func (r *PollRepository) PickPollQuestions(ctx context.Context) (pollQuestions []models.PollQuestion, err error) {
	funcName := "PickPollQuestions"
	query := `
	SELECT cq.question_id, cq.question_text, cq.type
	FROM csat_question AS cq;
	`

	rows, err := r.db.Query(ctx, query)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("PickPollQuestions (query): %w", err)
	}

	for rows.Next() {
		pollQuestion := models.PollQuestion{}
		if err := rows.Scan(&pollQuestion.QuestionID, &pollQuestion.QuestionText, &pollQuestion.QuestionType); err != nil {
			return nil, fmt.Errorf("PickPollQuestions (scan): %w", err)
		}
		pollQuestions = append(pollQuestions, pollQuestion)
	}

	return pollQuestions, nil
}

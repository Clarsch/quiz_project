package frontdto

import "quizzy_game/internal/dataTypes"

type AnswerResultDTO struct {
	QuizId        string                            `json:"quizId"`
	QuestionId    string                            `json:"questionId"`
	CorrectAnswer string                            `json:"correctAnswer"`
	ScoreBoard    map[string]*dataTypes.Participant `json:"scores"`
}

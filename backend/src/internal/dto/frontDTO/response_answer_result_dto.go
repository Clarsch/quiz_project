package frontdto

import (
	"encoding/json"
	"fmt"
	"quizzy_game/internal/dataTypes"
)

type AnswerResultDTO struct {
	QuizId        string                            `json:"quizId"`
	QuestionId    string                            `json:"questionId"`
	CorrectAnswer string                            `json:"correctAnswer"`
	UserScore     int                               `json:"userScore"`
	ScoreBoard    map[string]*dataTypes.Participant `json:"scoreboard"`
}

func (q AnswerResultDTO) JsonString() string {
	qJSON, err := json.MarshalIndent(q, "", "    ")
	if err != nil {
		return fmt.Sprintf("Error converting AnswerResultDTO to JSON: %v", err)
	}
	return string(qJSON)
}

package dbdto

import (
	"encoding/json"
	"fmt"
	"quizzy_game/internal/enums"
)

type QuestionIncomingDTO struct {
	Type          enums.QuestionType `json:"type"`
	Difficulty    string             `json:"difficulty"`
	Category      string             `json:"category"`
	Question      string             `json:"question"`
	CorrectAnswer string             `json:"correct_answer"`
	WrongAnswer   []string           `json:"incorrect_answers"`
}

func (q QuestionIncomingDTO) String() string {
	// Convert the Question struct to a JSON string for better visibility
	qJSON, err := json.MarshalIndent(q, "", "    ")
	if err != nil {
		return fmt.Sprintf("Error converting Question to JSON: %v", err)
	}
	return string(qJSON)
}

package frontdto

import (
	"encoding/json"
	"fmt"
	dbdto "quizzy_game/internal/dto/dbDTO"
	"quizzy_game/internal/enums"
)

type QuizOptionsDTO struct {
	Difficulty []enums.Difficulty          `json:"difficulties"`
	Type       []enums.QuestionType        `json:"types"`
	Category   []dbdto.CategoryIncomingDTO `json:"categories"`
}

func (q QuizOptionsDTO) String() string {
	qJSON, err := json.MarshalIndent(q, "", "    ")
	if err != nil {
		return fmt.Sprintf("Error converting QuizOptionsDTO to JSON: %v", err)
	}
	return string(qJSON)
}

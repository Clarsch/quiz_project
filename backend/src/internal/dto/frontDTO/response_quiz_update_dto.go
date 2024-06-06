package frontdto

import (
	"quizzy_game/internal/dataTypes"
	"quizzy_game/internal/enums/quizStatus"
)

type QuizUpdateResponseDTO struct {
	QuizId       string                            `json:"quizId"`
	Name         string                            `json:"name"`
	QuizStatus   quizStatus.QuizStatus             `json:"quizStatus"`
	Participants map[string]*dataTypes.Participant `json:"participants"`
}

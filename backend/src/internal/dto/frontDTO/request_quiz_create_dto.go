package frontdto

import (
	"quizzy_game/internal/enums/questionType"
	"quizzy_game/internal/enums/quizDifficulty"
)

type QuizCreateRequestDTO struct {
	Name       string                    `json:"name"`
	CategoryId int                       `json:"categoryId"`
	Difficulty quizDifficulty.Difficulty `json:"difficulty"`
	Type       questionType.QuestionType `json:"type"`
}

func (r QuizCreateRequestDTO) GetExample() string {
	return `
		{
			"request": "create",
			"data": {
				"name": "quiz1",
				"categoryId": 9,
				"difficulty": "easy",
				"type": "multiple"
			}
		}`

}

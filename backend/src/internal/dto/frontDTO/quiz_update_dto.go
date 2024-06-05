package frontdto

import (
	"quizzy_game/internal/dataTypes"
	"quizzy_game/internal/enums"
)

type QuizCreateRequestDTO struct {
	Name       string             `json:"name"`
	CategoryId string             `json:"categoryId"`
	Difficulty enums.Difficulty   `json:"difficulty"`
	Type       enums.QuestionType `json:"type"`
}

type QuizUpdateRequestDTO struct {
	Id      string `json:"quizId"`
	Command string `json:"command"`
}

type QuizUpdateResponseDTO struct {
	QuizId       string                            `json:"quizId"`
	Name         string                            `json:"name"`
	QuizStatus   enums.QuizStatus                  `json:"quizStatus"`
	Participants map[string]*dataTypes.Participant `json:"participants"`
}

type QuestionAnswerDTO struct {
	QuizId     string `json:"quizId"`
	QuestionId string `json:"questionId"`
	AnswerId   string `json:"answerId"`
}

type QuestionAnswerResultDTO struct {
	QuizId        string                            `json:"quizId"`
	QuestionId    string                            `json:"questionId"`
	CorrectAnswer string                            `json:"correctAnswer"`
	ScoreBoard    map[string]*dataTypes.Participant `json:"scores"`
}

type QuestionDTO struct {
	QuizId     string   `json:"quizId"`
	QuestionId string   `json:"questionId"`
	Question   string   `json:"question"`
	Options    []Option `json:"options"`
	// CorrectOptionId string   `json:"correctOptionId"`
}

type Option struct {
	OptionId string `json:"optionId"`
	Option   string `json:"option"`
}

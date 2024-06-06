package frontdto

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

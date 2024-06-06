package frontdto

type QuestionAnswerDTO struct {
	QuizId     string `json:"quizId"`
	QuestionId string `json:"questionId"`
	AnswerId   string `json:"answerId"`
}

func (r QuestionAnswerDTO) GetExample() string {
	return `
		{
			"request": "answer",
			"data": {
				"quizId": "4b78a3f7-0c31-437b-84a1-2a3b68724335",
				"questionId": "4b78a3f7-0c31-437b-84a1-2a3b68452315",
				"answerId": "4b78a3f7-0c31-437b-84a1-feaeek3d12312"
			}
		}`

}

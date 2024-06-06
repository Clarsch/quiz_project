package frontdto

type QuizUpdateRequestDTO struct {
	QuizId string `json:"quizId"`
}

func (r QuizUpdateRequestDTO) GetExample() string {
	return `
		{
			"request": "start",
			"data": {
				"quizId": "4b78a3f7-0c31-437b-84a1-2a3b68724335"
			}
		}`

}

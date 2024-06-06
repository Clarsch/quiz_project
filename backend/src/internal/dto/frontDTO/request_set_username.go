package frontdto

type SetUsernameRequestDTO struct {
	Username string `json:"username"`
}

func (r SetUsernameRequestDTO) GetExample() string {
	return `
		{
			"request": "set_username",
			"data": {
				"username": "Peter"
			}
		}`

}

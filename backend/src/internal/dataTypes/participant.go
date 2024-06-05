package dataTypes

type Participant struct {
	Ref   *User   `json:"-"`
	Name  *string `json:"name"`
	Score int     `json:"score"`
}

func NewParticipant(user *User) *Participant {
	return &Participant{user, &user.Name, 0}
}

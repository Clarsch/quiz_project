package dataTypes

import (
	"math/rand"
	dbdto "quizzy_game/internal/dto/dbDTO"
	"time"

	"github.com/google/uuid"
)

type Question struct {
	Id            string
	ref           *dbdto.QuestionIncomingDTO
	IsAskedStatus bool
	LastAskedTime time.Time
}

func NewQuestion(q dbdto.QuestionIncomingDTO) Question {
	questionId := uuid.NewString()
	return Question{questionId, &q, false, time.Time{}}
}

func (q Question) GetOptions() []string {
	options := append(q.ref.WrongAnswer, q.ref.CorrectAnswer)
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})
	return options
}

func (q Question) GetQuestion() string {
	return q.ref.Question
}

func (q Question) GetCorrectAnswer() string {
	return q.ref.CorrectAnswer
}

func (q Question) IsNotAsked() bool {
	return !q.IsAskedStatus
}

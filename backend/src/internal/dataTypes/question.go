package dataTypes

import (
	"math/rand"
	dbdto "quizzy_game/internal/dto/dbDTO"
	"time"

	"github.com/google/uuid"
)

type Question struct {
	Id              string
	ref             *dbdto.QuestionIncomingDTO
	Options         []AnswerOption
	CorrectOptionId string
	IsAskedStatus   bool
	LastAskedTime   time.Time
}

type AnswerOption struct {
	OptionId string
	Answer   string
}

func NewQuestion(q dbdto.QuestionIncomingDTO) Question {
	questionId := uuid.NewString()
	var options []AnswerOption

	correctOptionId := uuid.NewString()
	options = append(options, AnswerOption{correctOptionId, q.CorrectAnswer})
	for _, o := range q.WrongAnswer {
		optionId := uuid.NewString()
		options = append(options, AnswerOption{optionId, o})
	}
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return Question{
		questionId,
		&q,
		options,
		correctOptionId,
		false,
		time.Time{},
	}

}

func (q Question) GetQuestion() string {
	return q.ref.Question
}

func (q Question) GetCorrectAnswer() string {
	for _, o := range q.Options {
		if o.OptionId == q.CorrectOptionId {
			return o.Answer
		}
	}
	return "Could not find correct Answer"

}

func (q Question) IsNotAsked() bool {
	return !q.IsAskedStatus
}

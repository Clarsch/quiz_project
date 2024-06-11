package dataTypes

import (
	"encoding/json"
	"fmt"
	"math/rand"
	dbdto "quizzy_game/internal/dto/dbDTO"
	"time"

	"github.com/google/uuid"
)

type Question struct {
	Id              string                     `json:"questionId"`
	ref             *dbdto.QuestionIncomingDTO `json:"-"`
	QuestionString  *string                    `json:"question"`
	Options         []AnswerOption             `json:"options"`
	CorrectOptionId string                     `json:"correctOptionId"`
	IsAskedStatus   bool                       `json:"-"`
	LastAskedTime   time.Time                  `json:"-"`
}

type AnswerOption struct {
	OptionId string `json:"optionId"`
	Answer   string `json:"option"`
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
		&q.Question,
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

func (q Question) JsonString() string {
	qJSON, err := json.MarshalIndent(q, "", "    ")
	if err != nil {
		return fmt.Sprintf("Error converting Question to JSON: %v", err)
	}
	return string(qJSON)
}

func (q QuizQuestion) JsonString() string {
	qJSON, err := json.MarshalIndent(q, "", "    ")
	if err != nil {
		return fmt.Sprintf("Error converting QuizQuestion to JSON: %v", err)
	}
	return string(qJSON)
}

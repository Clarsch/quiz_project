package dataTypes

import (
	"encoding/json"
	"fmt"
	dbdto "quizzy_game/internal/dto/dbDTO"
	"quizzy_game/internal/enums"
	"strings"
)

type Quiz struct {
	Id            string                    `json:"id"`
	Name          string                    `json:"name"`
	QuizStatus    enums.QuizStatus          `json:"quizStatus"`
	Category      dbdto.CategoryIncomingDTO `json:"category"`
	Difficulty    enums.Difficulty          `json:"difficulty"`
	Type          enums.QuestionType        `json:"type"`
	Questions     map[string]*Question      `json:"questions"`
	Participants  map[string]*Participant   `json:"participants"`
	StatusChannel *chan enums.QuizStatus    `json:"-"`
}

func (q Quiz) String() string {
	qJSON, err := json.MarshalIndent(q, "", "    ")
	if err != nil {
		return fmt.Sprintf("Error converting Quiz to JSON: %v", err)
	}
	return string(qJSON)
}

func (q Quiz) RemainingQuestions() int {
	unaskedCounter := 0
	for _, question := range q.Questions {
		if question.IsNotAsked() {
			unaskedCounter++
		}
	}

	return unaskedCounter
}

func (q Quiz) ParticipantsAsString() string {
	var names []string
	for _, pt := range q.Participants {
		names = append(names, pt.Ref.Name)
	}
	return strings.Join(names, ", ")
}

func (q Quiz) ScoreBoard() []Participant {
	var participants []Participant
	for _, p := range q.Participants {
		participants = append(participants, *p)
	}
	return participants

	// var scoresList []string
	// for _, pt := range q.Participants {
	// 	scoresList = append(scoresList, fmt.Sprintf("%s: %d", pt.Ref.Name, pt.Score))
	// }
	// return strings.Join(scoresList, ", ")

}

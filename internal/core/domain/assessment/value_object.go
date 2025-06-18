package assessment

import (
	"errors"
	"fmt"
)

type Id int
type Content string

type QuestionType string

type NoOfQuestions int

func (n NoOfQuestions) isEmpty() error {
	if n == 0 {
		return errors.New("number of questions has to be greater than 0")
	}
	return nil

}
func (n NoOfQuestions) isMax(amount int) (*NoOfQuestions, error) {
	if int(n) > amount {
		return nil, fmt.Errorf("number of questions has to be less than %d", amount)
	}
	return &n, nil

}
func (n NoOfQuestions) isMin(amount int) (*NoOfQuestions, error) {
	if int(n) < amount {
		return nil, fmt.Errorf("number of questions has to be greater than %d", amount)
	}
	return &n, nil

}

var (
	Essay                 QuestionType = "essay"
	TrueFalse             QuestionType = "true/false"
	MultiAnswer           QuestionType = "checkbox"
	OneAnswer             QuestionType = "radio"
	FillInTheBlank        QuestionType = "fill-in-the-blank"
	MatchQuestionToOption QuestionType = "match-questions-to-options"
)

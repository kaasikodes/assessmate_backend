package assessment

import (
	"errors"
	"fmt"
)

type QuestionValidator struct {
	value NoOfQuestions
	err   error
}

func NewQuestionValidator(value NoOfQuestions) *QuestionValidator {
	return &QuestionValidator{value: value}
}

func (q *QuestionValidator) IsEmpty() *QuestionValidator {
	if q.err != nil {
		return q
	}
	if q.value == 0 {
		q.err = errors.New("number of questions has to be greater than 0")
	}
	return q
}

func (q *QuestionValidator) IsMax(max int) *QuestionValidator {
	if q.err != nil {
		return q
	}
	if int(q.value) > max {
		q.err = fmt.Errorf("number of questions has to be less than or equal to %d", max)
	}
	return q
}

func (q *QuestionValidator) IsMin(min int) *QuestionValidator {
	if q.err != nil {
		return q
	}
	if int(q.value) < min {
		q.err = fmt.Errorf("number of questions has to be greater than or equal to %d", min)
	}
	return q
}

func (q *QuestionValidator) Error() error {
	return q.err
}

func (q *QuestionValidator) Value() NoOfQuestions {
	return q.value
}

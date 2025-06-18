package assessment

type Option struct {
	id        Id
	content   Content
	isCorrect bool
}
type OneAnswerQuestion struct {
	id      Id
	content Content

	option Option
}
type MultiAnswerQuestion struct {
	id      Id
	content Content

	options []Option
}
type TrueFalseQuestion struct {
	id      Id
	content Content

	options []Option
}

type EssayQuestion struct {
	id              Id
	content         Content
	suggestedAnswer Content
}

func NewEssayQuestion(content, suggestedAnswer Content) EssayQuestion {
	// validate the data input say content.isMinLength
	// suggestAnswer.isEmpty
	//TODO: Pile up errs using validation_errs and return as err, also test to ensure the errors can be decoded and sent to client as an array of errors
	return EssayQuestion{
		content:         content,
		suggestedAnswer: suggestedAnswer,
	}

}

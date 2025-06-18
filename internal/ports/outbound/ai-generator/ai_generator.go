package aigenerator

import (
	"context"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/assessment"
)

// TODO: Implement adapter for this that will take in the server url (internal), also reearch if ollama allows for grpc calls to speed things up
type AiGenerator interface {
	GenerateEssayQuestions(ctx context.Context, noOfQuestions assessment.NoOfQuestions) ([]assessment.EssayQuestion, error)
	GenerateMultiAnswerQuestions(ctx context.Context, noOfQuestions assessment.NoOfQuestions) ([]assessment.MultiAnswerQuestion, error)
	GenerateOneAnswerQuestions(ctx context.Context, noOfQuestions assessment.NoOfQuestions) ([]assessment.OneAnswerQuestion, error)
	GenerateTrueFalseQuestions(ctx context.Context, noOfQuestions assessment.NoOfQuestions) ([]assessment.TrueFalseQuestion, error)
}

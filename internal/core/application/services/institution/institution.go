package institution

import (
	"context"
	"fmt"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/institution"
	institute_repo "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/institution"
)

// DTOs (used as input/output to/from service methods)

type CreateInstitutionRequest struct {
	Name        string
	Description string
	Email       string
}

type CreateInstitutionResponse struct {
	Id          int
	Name        string
	Description string
	Email       string
	CreatedAt   string
	UpdatedAt   string
}

type InstitutionManagementService struct {
	instituteRepo institute_repo.InstitutionRepository
}

// Constructor
func NewInstitutionManagementService(repo institute_repo.InstitutionRepository) *InstitutionManagementService {
	return &InstitutionManagementService{
		instituteRepo: repo,
	}
}

// CreateInstitution handles the creation of a new institution.
func (s *InstitutionManagementService) CreateInstitution(ctx context.Context, req CreateInstitutionRequest) (*CreateInstitutionResponse, error) {
	name, err := institution.NewName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid institution name: %w", err)
	}

	desc, err := institution.NewDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("invalid description: %w", err)
	}

	email, err := institution.NewEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	inst, err := institution.NewInstitution(name, desc, email)
	if err != nil {
		return nil, err
	}

	created, err := s.instituteRepo.CreateInstitution(ctx, inst)
	if err != nil {
		return nil, fmt.Errorf("failed to create institution: %w", err)
	}

	return &CreateInstitutionResponse{
		Id:          created.Id().Value(),
		Name:        created.Name().String(),
		Description: created.Description().String(),
		Email:       created.Email().String(),
		CreatedAt:   created.CreatedAt().String(),
		UpdatedAt:   created.UpdatedAt().String(),
	}, nil
}

// create institution
// invite staff to institution
// create staff in institution
// blacklist users in instititution
// activate users in instituion
// change user status in institution
// audit trail of users in institution
// lecture material creates by users in institution
// create categories in institution
// create groups in institution
// group users in institution (users in group can only create assessments in this category)
// be able to a list staff profiles within the institution, assign to groups ...
// view assessments within the institution(be able to filter by categoryId, staffId(userId))

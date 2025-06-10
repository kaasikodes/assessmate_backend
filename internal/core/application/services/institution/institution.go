package institution

import (
	"context"
	"fmt"
	"time"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/institution"
	"github.com/kaasikodes/assessmate_backend/internal/ports/outbound/email"
	institute_repo "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/institution"
	errors "github.com/kaasikodes/assessmate_backend/internal/shared"
)

// DTOs (used as input/output to/from service methods)

type (
	InviteStaffRequest struct {
		Staff         []StaffInvite
		InstitutionId int
	}
	StaffInvite struct {
		Name, Email string
	}

	CreateInstitutionRequest struct {
		Name, Email, Description string
	}
	AddStaffToGroupRequest struct {
		InstitutionId, GroupId, StaffId int
	}
	CreateGroupRequest struct {
		Name, Description string
		InstitutionId     int
	}

	CreateGroupResponse struct {
		Id          int
		Name        string
		Description string
	}
	CreateInstitutionResponse struct {
		Id          int
		Name        string
		Description string
		Email       string
		CreatedAt   string
		UpdatedAt   string
	}
	Staff struct {
		Id        int
		Name      string
		Email     string
		Status    string
		CreatedAt string
		UpdatedAt string
	}
	Group struct {
		Id        int
		Name      string
		CreatedAt string
		UpdatedAt string
	}
	StaffResponse struct {
		Staff []Staff
		Total int
	}
	GroupResponse struct {
		Groups []Group
		Total  int
	}
)

type InstitutionManagementService struct {
	instituteRepo institute_repo.InstitutionRepository
	emailClient   email.EmailClient
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

// get InstitutionById
func (s *InstitutionManagementService) GetInstitutionById(ctx context.Context, id int) (*CreateInstitutionResponse, error) {
	instituteId, err := institution.NewId(id)
	if err != nil {
		return nil, err
	}
	data, err := s.instituteRepo.GetInstitutionById(ctx, instituteId)
	if err != nil {
		return nil, err
	}
	return &CreateInstitutionResponse{
		Id:          data.Id().Value(),
		Name:        data.Name().String(),
		Description: data.Description().String(),
		Email:       data.Email().String(),
		CreatedAt:   data.CreatedAt().Format(time.RFC1123),
		UpdatedAt:   data.UpdatedAt().Format(time.RFC1123),
	}, nil

}

// get staff
// TODO Add ability to filter staff by status
func (s *InstitutionManagementService) GetStaff(ctx context.Context, id int) (*StaffResponse, error) {

	instituteId, err := institution.NewId(id)
	if err != nil {
		return nil, nil
	}
	data, total, err := s.instituteRepo.ListStaff(ctx, instituteId)
	if err != nil {
		return nil, err
	}

	staff := make([]Staff, len(data))
	for i, d := range data {
		staff[i].Id = d.Id().Value()
		staff[i].Name = d.Name().String()
		staff[i].Email = d.Email().String()
		staff[i].Status = d.Status().String()
		staff[i].CreatedAt = d.CreatedAt().String()
		staff[i].UpdatedAt = d.UpdatedAt().String()

	}
	return &StaffResponse{
		Total: total,
		Staff: staff,
	}, nil

}

// create group
func (s *InstitutionManagementService) CreateGroup(ctx context.Context, req CreateGroupRequest) (*CreateGroupResponse, error) {
	var valErrs errors.ValidationErrors

	instituteId, err := institution.NewId(req.InstitutionId)
	if err != nil {
		valErrs.Add("instituteId", err.Error())

	}

	name, err := institution.NewName(req.Name)
	if err != nil {
		valErrs.Add("name", err.Error())

	}

	desc, err := institution.NewDescription(req.Description)
	if err != nil {
		valErrs.Add("description", err.Error())
	}

	group, err := institution.NewGroup(name, desc)
	if valErrs.HasErrors() {
		return nil, &valErrs

	}

	created, err := s.instituteRepo.CreateGroup(ctx, instituteId, group)
	if err != nil {
		return nil, fmt.Errorf("failed to create institution: %w", err)
	}

	return &CreateGroupResponse{
		Id:          created.Id().Value(),
		Name:        created.Name().String(),
		Description: created.Description().String(),
	}, nil
}

// add staff to group
func (s *InstitutionManagementService) AddStaffToGroup(ctx context.Context, req AddStaffToGroupRequest) error {
	var valErrs errors.ValidationErrors

	instituteId, err := institution.NewId(req.InstitutionId)
	if err != nil {
		valErrs.Add("instituteId", err.Error())

	}
	groupId, err := institution.NewId(req.GroupId)
	if err != nil {
		valErrs.Add("groupId", err.Error())

	}
	staffId, err := institution.NewId(req.StaffId)
	if err != nil {
		valErrs.Add("staffId", err.Error())

	}

	if valErrs.HasErrors() {
		return &valErrs

	}

	err = s.instituteRepo.AddStaffToGroup(ctx, instituteId, groupId, staffId)
	if err != nil {
		return fmt.Errorf("failed to add staff to group: %w", err)
	}
	return nil
}

// list groups
func (s *InstitutionManagementService) GetGroups(ctx context.Context, id int) (*GroupResponse, error) {

	instituteId, err := institution.NewId(id)
	if err != nil {
		return nil, nil
	}
	data, total, err := s.instituteRepo.ListGroups(ctx, instituteId)
	if err != nil {
		return nil, err
	}

	groups := make([]Group, len(data))
	for i, d := range data {
		groups[i].Id = d.Id().Value()
		groups[i].Name = d.Name().String()
		groups[i].CreatedAt = d.CreatedAt().String()
		groups[i].UpdatedAt = d.UpdatedAt().String()

	}
	return &GroupResponse{
		Total:  total,
		Groups: groups,
	}, nil

}

//TODO: 3 add accessible course to group
//TODO: 2 list accessible courses in group

//TODO: 1 invite staff
// func (s *InstitutionManagementService) InviteStaff(ctx context.Context, req CreateInstitutionRequest) error {
// 	var err errors.ValidationErrors

// 	//  will recieve and run the email sending as a background job

// }

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

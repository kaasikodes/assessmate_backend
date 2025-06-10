package institution

import (
	"context"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/institution"
)

// InstitutionRepository defines the contract for interacting with institution aggregates.
type InstitutionRepository interface {
	// Institution
	CreateInstitution(ctx context.Context, payload *institution.Institution) (*institution.Institution, error)
	GetInstitutionById(ctx context.Context, id institution.Id) (*institution.Institution, error)
	UpdateInstitution(ctx context.Context, inst *institution.Institution) error
	DeleteInstitution(ctx context.Context, id institution.Id) error

	// Staff
	AddStaffToInstitution(ctx context.Context, institutionId institution.Id, staff institution.Staff) error
	GetStaffById(ctx context.Context, institutionId, staffId institution.Id) (*institution.Staff, error)
	RemoveStaffFromInstitution(ctx context.Context, institutionId, staffId institution.Id) error
	ListStaff(ctx context.Context, institutionId institution.Id) ([]institution.Staff, int, error)

	// Group
	CreateGroup(ctx context.Context, institutionId institution.Id, group *institution.Group) (*institution.Group, error)
	GetGroupById(ctx context.Context, institutionId, groupId institution.Id) (*institution.Group, error)
	UpdateGroup(ctx context.Context, institutionId institution.Id, group *institution.Group) error
	DeleteGroup(ctx context.Context, institutionId, groupId institution.Id) error
	AddStaffToGroup(ctx context.Context, institutionId, groupId, staffId institution.Id) error
	RemoveStaffFromGroup(ctx context.Context, institutionId, groupId, staffId institution.Id) error
	ListGroups(ctx context.Context, institutionId institution.Id) ([]institution.Group, int, error)

	// Course Access (Group <-> Course)
	AddAccessibleCourseToGroup(ctx context.Context, institutionId, groupId, courseId institution.Id) error
	RemoveAccessibleCourseFromGroup(ctx context.Context, institutionId, groupId, courseId institution.Id) error
	ListAccessibleCoursesForGroup(ctx context.Context, institutionId, groupId institution.Id) ([]institution.Course, error)
}

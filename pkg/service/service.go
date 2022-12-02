package service

import (
	"github.com/RucardTomsk/course_book/model"
	"github.com/RucardTomsk/course_book/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) error
	GenerateToken(email string, password string) (string, string, error)
	ParseToken(accessToken string) (string, error)
	GetUserFioByGuid(guid string) (string, error)
	GetUserNotAccess(guid_node string) ([]model.User, error)
	GenerateRefreshToken() (string, error)
	IssueSessionUser(user model.User, refreshToken string) error
	GetUserToRefreshToken(refreshToken string) (model.User, error)
	CreateResetPassword(user model.User) error
	CheckResetPassword(code string, user model.User) error
	UserResetPassword(user model.User, newPassword string) error
	GetUserByEmail(email string) (model.User, error)
}

type Plans interface {
	CreatePlans(NameDiscipline string, ByteTable []byte, guid_faculty string) error
	GetPlans(guid_programm string) (map[string][]model.BriefPlan, error)
	GetWorkProgram(guid_plan string) (model.FullPlan, error)
	SavePlan(guid_plan string, key_field string, text string) error
	GetField(guid_plan string, key_field string) (string, error)
	GetNamePlans(guid string) ([]string, error)
	CloneFieldPlan(guid_from, guid_to string) error
}

type Faculty interface {
	GetMasFaculty() ([]model.Faculty, error)
	GetNameFaculty(guid string) (string, error)
}

type Program interface {
	GetMasProgram(guid_faculty string) ([]model.Program, error)
	GetNameProgramAndFaculty(guid_program string) ([]string, error)
}

type Role interface {
	CheckRoleAdmin(guid_user string) (bool, error)
	IssueAccess(guid_user, guid_node string) (string, error)
	CheckAccess(guid_user, guid_node string) (bool, error)
	CreateInvite(guid_node string) (string, error)
	UseInvite(guid_invite, guid_user string) error
}
type Service struct {
	Authorization
	Plans
	Faculty
	Program
	Role
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Plans:         NewPlansService(repos.Plans),
		Faculty:       NewFaculteService(repos.Faculty),
		Program:       NewProgramService(repos.Program),
		Role:          NewRoleService(repos.Role),
	}
}

package repository

import (
	"github.com/RucardTomsk/course_book/model"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Authorization interface {
	CreateUser(user model.User) error
	GetUser(username string, password string) (model.User, error)
	GetUserFIOByGuid(guid string) (string, error)
	CheckAbsentEmail(email string) (bool, error)
	GetUserNotAccess(guid_node string) ([]model.User, error)
	IssueSessionUser(user model.User, refreshToken string) error
	GetUserToRefreshToken(refreshToken string) (model.User, error)
}

type Plans interface {
	CreatePlan(dict map[string]string) (string, error)
	CreateRelationship(guid_node_a string, guid_node_b string, typeR string) error
	CreateProgramm(name string, directions string) (string, error)
	GetMasPlan(guid_program string) ([]model.BriefPlan, error)
	GetWorkProgram(guid_plan string) (model.FullPlan, error)
	SavePlan(guid_plan string, key_field string, text string) error
	GetField(guid_plan string, key_field string) (string, error)
	GetNamePlans(guid string) ([]string, error)
}

type Faculty interface {
	GetMasFaculte() ([]model.Faculty, error)
	GetNameFaculte(guid string) (string, error)
}

type Program interface {
	GetMasProgram(guid_faculty string) ([]model.Program, error)
	GetNameProgramAndFaculty(guid_program string) ([]string, error)
}

type Role interface {
	IssueAccess(guid_user, guid_node string) (string, error)
	CheckRoleAdmin(guid_user string) (bool, error)
	CheckAccess(guid_user, guid_node string) (bool, error)
}

type Repository struct {
	Authorization
	Plans
	Faculty
	Program
	Role
}

func NewRepository(driver *neo4j.Driver) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(driver),
		Plans:         NewPlansRepository(driver),
		Faculty:       NewFacultyRepository(driver),
		Program:       NewProgramRepository(driver),
		Role:          NewRoleRepository(driver),
	}
}

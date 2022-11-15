package service

import (
	"github.com/RucardTomsk/course_book/model"
	"github.com/RucardTomsk/course_book/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Plans interface {
	CreatePlans(NameDiscipline string, ByteTable []byte) error
	GetPlans(guid_programm string) ([]model.BriefPlan, error)
	GetWorkProgram(guid_plan string) (model.FullPlan, error)
	SavePlan(guid_plan string, key_field string, text string) error
	GetField(guid_plan string, key_field string) (string, error)
}

type Service struct {
	Authorization
	Plans
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Plans:         NewPlansService(repos.Plans),
	}
}

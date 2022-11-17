package service

import (
	"github.com/RucardTomsk/course_book/model"
	"github.com/RucardTomsk/course_book/pkg/repository"
)

type ProgramService struct {
	repo repository.Program
}

func NewProgramService(repo repository.Program) *ProgramService {
	return &ProgramService{repo: repo}
}

func (s *ProgramService) GetMasProgram(guid_faculty string) ([]model.Program, error) {
	return s.repo.GetMasProgram(guid_faculty)
}

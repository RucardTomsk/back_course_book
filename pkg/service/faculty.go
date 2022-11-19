package service

import (
	"github.com/RucardTomsk/course_book/model"
	"github.com/RucardTomsk/course_book/pkg/repository"
)

type FacultyService struct {
	repo repository.Faculty
}

func NewFaculteService(repo repository.Faculty) *FacultyService {
	return &FacultyService{repo: repo}
}

func (s *FacultyService) GetMasFaculty() ([]model.Faculty, error) {
	return s.repo.GetMasFaculte()
}

func (s *FacultyService) GetNameFaculty(guid string) (string, error) {
	return s.repo.GetNameFaculte(guid)
}

package service

import "github.com/RucardTomsk/course_book/pkg/repository"

type RoleService struct {
	repo repository.Role
}

func NewRoleService(repo repository.Role) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) CheckRoleAdmin(guid_user string) (bool, error) {
	return s.repo.CheckRoleAdmin(guid_user)
}

func (s *RoleService) IssueAccess(guid_user, guid_node string) (string, error) {
	return s.repo.IssueAccess(guid_user, guid_node)
}

func (s *RoleService) CheckAccess(guid_user, guid_node string) (bool, error) {
	return s.repo.CheckAccess(guid_user, guid_node)
}

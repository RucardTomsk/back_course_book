package repository

import (
	"github.com/RucardTomsk/course_book/model"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type FacultyRepository struct {
	driver *neo4j.Driver
}

func NewFacultyRepository(driver *neo4j.Driver) *FacultyRepository {
	return &FacultyRepository{driver: driver}
}

func (r *FacultyRepository) GetMasFaculte() ([]model.Faculty, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	var list []model.Faculty
	result, err := session.Run("MATCH (f:Faculty) RETURN f", map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	for result.Next() {
		prop := result.Record().Values[0].(neo4j.Node).Props

		list = append(list, model.Faculty{
			Name: prop["Name"].(string),
			Guid: prop["guid"].(string),
			Icon: prop["icon"].(string),
		})
	}

	return list, nil
}

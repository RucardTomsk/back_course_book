package repository

import (
	"github.com/RucardTomsk/course_book/model"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type ProgramRepository struct {
	driver *neo4j.Driver
}

func NewProgramRepository(driver *neo4j.Driver) *ProgramRepository {
	return &ProgramRepository{driver: driver}
}

func (r *ProgramRepository) GetMasProgram(guid_faculty string) ([]model.Program, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	var list []model.Program
	result, err := session.Run("MATCH (program)-[]->(faculty) WHERE faculty.guid = $guid RETURN program", map[string]interface{}{
		"guid": guid_faculty,
	})
	if err != nil {
		return nil, err
	}

	for result.Next() {
		prop := result.Record().Values[0].(neo4j.Node).Props
		list = append(list, model.Program{
			Name:       prop["Name"].(string),
			Directions: prop["Directions"].(string),
			Guid:       prop["guid"].(string),
		})
	}

	return list, nil
}

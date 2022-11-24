package repository

import (
	"fmt"
	"unicode"

	"github.com/RucardTomsk/course_book/model"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type PlansRepository struct {
	driver *neo4j.Driver
}

func NewPlansRepository(driver *neo4j.Driver) *PlansRepository {
	return &PlansRepository{driver: driver}
}

func (r *PlansRepository) CreatePlan(dict map[string]string) (string, error) {
	session := GetSession(*r.driver)
	defer session.Close()
	mas_key := []string{"DirectionTraining", "TrainingProfile", "FormTraining", "Qualification", "Year", "Name", "PurposeMastering", "ResultsMastering", "PlaceDiscipline", "SemesterMastering", "Code", "ImplementationLanguage", "DevelopmentMastering", "EntranceRequirements", "ScopeDiscipline", "ContentDiscipline", "CurrentControl", "EvaluationProcedure", "MethodologicalSupport", "References", "ListInformationTechnologies", "MaterialSupport", "InformationDevelopers", "NameFaculty"}
	guid := uuid.New().String()
	map_key := make(map[string]interface{})
	for _, key := range mas_key {
		if key == "ScopeDiscipline" {
			map_key[key] = fmt.Sprintf("Общая трудоемкость дисциплины составляет %s з.е., %s часов, из которых:\n"+
				"– лекции: %s ч.;\n"+
				"– семинарские занятия: %s ч.\n"+
				"– практические занятия: %s ч.\n"+
				"– лабораторные работы: %s ч.\n"+
				"в том числе практическая подготовка: 0 ч.\n"+
				"Объем самостоятельной работы студента определен учебным планом.",
				dict["ZE"],
				dict["K"],
				dict["L"],
				dict["C"],
				dict["P"],
				dict["LR"])
		} else {
			map_key[key] = dict[key]
		}
	}
	map_key["guid"] = guid
	_, err := session.Run("CREATE (n:Plan {guid: $guid, DirectionTraining: $DirectionTraining, TrainingProfile: $TrainingProfile, FormTraining: $FormTraining, Qualification: $Qualification, Year: $Year, Name: $Name, PurposeMastering: $PurposeMastering, ResultsMastering: $ResultsMastering, PlaceDiscipline: $PlaceDiscipline, SemesterMastering: $SemesterMastering, Code: $Code, ImplementationLanguage: $ImplementationLanguage, DevelopmentMastering: $DevelopmentMastering, EntranceRequirements: $EntranceRequirements, ScopeDiscipline: $ScopeDiscipline, ContentDiscipline: $ContentDiscipline, CurrentControl: $CurrentControl, EvaluationProcedure: $EvaluationProcedure, MethodologicalSupport: $MethodologicalSupport, References: $References, ListInformationTechnologies: $ListInformationTechnologies, MaterialSupport: $MaterialSupport, InformationDevelopers: $InformationDevelopers, NameFaculty: $NameFaculty})", map_key)
	if err != nil {
		return "", err
	}
	map_key["guid"] = guid + "_html"
	_, err = session.Run("CREATE (n:PlanHTML {guid: $guid, DirectionTraining: $DirectionTraining, TrainingProfile: $TrainingProfile, FormTraining: $FormTraining, Qualification: $Qualification, Year: $Year, Name: $Name, PurposeMastering: $PurposeMastering, ResultsMastering: $ResultsMastering, PlaceDiscipline: $PlaceDiscipline, SemesterMastering: $SemesterMastering, Code: $Code, ImplementationLanguage: $ImplementationLanguage, DevelopmentMastering: $DevelopmentMastering, EntranceRequirements: $EntranceRequirements, ScopeDiscipline: $ScopeDiscipline, ContentDiscipline: $ContentDiscipline, CurrentControl: $CurrentControl, EvaluationProcedure: $EvaluationProcedure, MethodologicalSupport: $MethodologicalSupport, References: $References, ListInformationTechnologies: $ListInformationTechnologies, MaterialSupport: $MaterialSupport, InformationDevelopers: $InformationDevelopers, NameFaculty: $NameFaculty})", map_key)
	if err != nil {
		return "", err
	}

	return guid, nil
}

func (r *PlansRepository) GetNamePlans(guid string) ([]string, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	result, err := session.Run("MATCH (p:Plan)-[]-(pr:Programm)-[]->(f:Faculty) WHERE p.guid = $guid RETURN f.Name,pr.Name,p.Name,f.guid", map[string]interface{}{
		"guid": guid,
	})
	if err != nil {
		return nil, err
	}

	var mas []string
	if result.Next() {
		mas = append(mas, result.Record().Values[0].(string))
		mas = append(mas, result.Record().Values[1].(string))
		mas = append(mas, result.Record().Values[2].(string))
		mas = append(mas, result.Record().Values[3].(string))
	}

	return mas, nil
}

func (r *PlansRepository) CreateProgramm(name string, directions string) (string, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	guid := uuid.New().String()
	_, err := session.Run("CREATE (:Programm {Name: $name, Directions: $directions, guid: $guid})", map[string]interface{}{
		"name":       name,
		"directions": directions,
		"guid":       guid,
	})

	if err != nil {
		return "", err
	}

	return guid, nil
}

func (r *PlansRepository) CreateRelationship(guid_node_a string, guid_node_b string, typeR string) error {
	session := GetSession(*r.driver)
	defer session.Close()

	_, err := session.Run(fmt.Sprintf("MATCH (a),(b) WHERE a.guid = $name_1 AND b.guid = $name_2 CREATE (a)-[:%s]->(b)", typeR), map[string]interface{}{
		"name_1": guid_node_a,
		"name_2": guid_node_b,
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *PlansRepository) GetMasPlan(guid_program string) ([]model.BriefPlan, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	var list []model.BriefPlan
	result, err := session.Run("MATCH (plan:Plan)-[]->(programm) WHERE programm.guid = $guid_programm RETURN plan", map[string]interface{}{
		"guid_programm": guid_program,
	})
	if err != nil {
		return nil, err
	}
	for result.Next() {
		prop := result.Record().Values[0].(neo4j.Node).Props

		list = append(list, model.BriefPlan{
			Code:              prop["Code"].(string),
			Name:              prop["Name"].(string),
			Guid:              prop["guid"].(string),
			SemesterMastering: prop["SemesterMastering"].(string),
		})
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *PlansRepository) GetWorkProgram(guid_plan string) (model.FullPlan, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	var workProgram model.FullPlan
	result, err := session.Run("MATCH (plan) WHERE plan.guid = $guid_plan RETURN plan", map[string]interface{}{
		"guid_plan": guid_plan,
	})
	if err != nil {
		return model.FullPlan{}, err
	}

	if result.Next() {
		props := result.Record().Values[0].(neo4j.Node).Props

		workProgram.Code = props["Code"].(string)
		workProgram.Name = props["Name"].(string)
		workProgram.DirectionTraining = props["DirectionTraining"].(string)
		workProgram.TrainingProfile = props["TrainingProfile"].(string)
		workProgram.FormTraining = props["FormTraining"].(string)
		workProgram.Qualification = props["Qualification"].(string)
		workProgram.Year = props["Year"].(string)
		workProgram.Guid = props["guid"].(string)
		workProgram.PurposeMastering = props["PurposeMastering"].(string)
		workProgram.ResultsMastering = props["ResultsMastering"].(string)
		workProgram.PlaceDiscipline = props["PlaceDiscipline"].(string)
		workProgram.SemesterMastering = props["SemesterMastering"].(string)
		workProgram.ImplementationLanguage = props["ImplementationLanguage"].(string)
		workProgram.DevelopmentMastering = props["DevelopmentMastering"].(string)
		workProgram.EntranceRequirements = props["EntranceRequirements"].(string)
		workProgram.ScopeDiscipline = props["ScopeDiscipline"].(string)
		workProgram.ContentDiscipline = props["ContentDiscipline"].(string)
		workProgram.CurrentControl = props["CurrentControl"].(string)
		workProgram.EvaluationProcedure = props["EvaluationProcedure"].(string)
		workProgram.MethodologicalSupport = props["MethodologicalSupport"].(string)
		workProgram.References = props["References"].(string)
		workProgram.ListInformationTechnologies = props["ListInformationTechnologies"].(string)
		workProgram.MaterialSupport = props["MaterialSupport"].(string)
		workProgram.InformationDevelopers = props["InformationDevelopers"].(string)
		workProgram.NameFaculty = props["NameFaculty"].(string)
	}

	return workProgram, nil
}

func (r *PlansRepository) SavePlan(guid_plan string, key_field string, text string) error {
	session := GetSession(*r.driver)
	defer session.Close()
	ru := []rune(key_field)
	ru[0] = unicode.ToUpper(ru[0])
	key_field = string(ru)
	_, err := session.Run(fmt.Sprintf("MATCH (plan) WHERE plan.guid = $guid_plan SET plan.%s = $text", key_field), map[string]interface{}{
		"guid_plan": guid_plan,
		"text":      text,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *PlansRepository) GetField(guid_plan string, key_field string) (string, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	var field string
	ru := []rune(key_field)
	ru[0] = unicode.ToUpper(ru[0])
	key_field = string(ru)
	result, err := session.Run(fmt.Sprintf("MATCH (plan) WHERE plan.guid = $guid_plan RETURN plan.%s", key_field), map[string]interface{}{
		"guid_plan": guid_plan,
	})
	if err != nil {
		return "", err
	}

	if result.Next() {
		field = result.Record().Values[0].(string)
	}

	return field, nil
}

package repository

import (
	"github.com/RucardTomsk/course_book/model"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type AuthRepository struct {
	driver *neo4j.Driver
}

func NewAuthRepository(driver *neo4j.Driver) *AuthRepository {
	return &AuthRepository{driver: driver}
}

func (r *AuthRepository) CreateUser(user model.User) error {
	session := GetSession(*r.driver)
	defer session.Close()
	guid := uuid.New().String()
	_, err := session.Run("CREATE (n:User {FIO: $1, Email: $3, Password: $4, guid: $5})", map[string]interface{}{
		"1": user.FIO,
		"3": user.Email,
		"4": user.Password,
		"5": guid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetUser(email string, password string) (model.User, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	var user model.User
	result, err := session.Run("MATCH (n:User) WHERE n.Email = $1 and n.Password = $2 RETURN n.guid, n.FIO", map[string]interface{}{
		"1": email,
		"2": password,
	})

	if err != nil {
		return user, err
	}

	if result.Next() {
		user.Guid = result.Record().Values[0].(string)
		user.FIO = result.Record().Values[1].(string)
		user.Email = email
		user.Password = password
	} else {
		return user, ErrRecordNotFound
	}

	return user, nil
}

func (r *AuthRepository) CheckAbsentEmail(email string) (bool, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	result, err := session.Run("MATCH (n:User) WHERE n.Email = $1 RETURN n.guid", map[string]interface{}{
		"1": email,
	})
	if err != nil {
		return false, err
	}
	if result.Next() {
		return false, nil
	} else {
		return true, nil
	}
}

func (r *AuthRepository) GetUserNotAccess(guid_node string) ([]model.User, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	var mas_user []model.User
	result, err := session.Run("MATCH (n) WHERE n.guid=$guid_node MATCH (u:User) WHERE NOT (u)-[:access]->(n) AND (:Admin)-[:role]->(u) RETURN u", map[string]interface{}{
		"guid_node": guid_node,
	})
	if err != nil {
		return nil, err
	}
	for result.Next() {
		props := result.Record().Values[0].(neo4j.Node).Props

		mas_user = append(mas_user, model.User{
			FIO:  props["FIO"].(string),
			Guid: props["guid"].(string),
		})
	}

	return mas_user, nil
}

func (r *AuthRepository) GetUserFIOByGuid(guid string) (string, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	result, err := session.Run("MATCH (n:User) WHERE n.guid = $1 RETURN n.FIO", map[string]interface{}{
		"1": guid,
	})

	if err != nil {
		return "", err
	}

	if result.Next() {
		return result.Record().Values[0].(string), nil
	} else {
		return "", ErrRecordNotFound
	}
}

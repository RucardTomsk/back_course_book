package repository

import (
	"github.com/RucardTomsk/course_book/model"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type AuthRepository struct {
	driver *neo4j.Driver
}

func NewAuthRepository(driver *neo4j.Driver) *AuthRepository {
	return &AuthRepository{driver: driver}
}

func (r *AuthRepository) CreateUser(user model.User) (int, error) {
	session := GetSession(*r.driver)
	defer session.Close()
	result, err := session.Run("CREATE (n:User {Name: $1, Username: $2,Password: $3}) RETURN id(n)", map[string]interface{}{
		"1": user.Name,
		"2": user.Username,
		"3": user.Password,
	})
	if err != nil {
		return 0, err
	}
	var id int
	if result.Next() {
		id = result.Record().Values[0].(int)
	} else {
		return 0, ErrRecordNotFound
	}
	user.Id = id

	return id, nil
}

func (r *AuthRepository) GetUser(username string, password string) (model.User, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	var user model.User
	result, err := session.Run("MATCH (n:User) WHERE n.Username = $1 and n.Password = $2 RETURN id(n), n.Name", map[string]interface{}{
		"1": username,
		"2": password,
	})

	if err != nil {
		return user, err
	}

	var id int
	if result.Next() {
		id = result.Record().Values[0].(int)
		user.Name = result.Record().Values[1].(string)
		user.Username = username
		user.Password = password
	} else {
		return user, ErrRecordNotFound
	}
	user.Id = id

	return user, nil
}

package repository

import (
	"fmt"
	"strings"

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

func (r *AuthRepository) GetUserByEmail(email string) (model.User, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	var user model.User
	result, err := session.Run("MATCH (n:User) WHERE n.Email = $1 RETURN n.guid, n.FIO", map[string]interface{}{
		"1": email,
	})

	if err != nil {
		return user, err
	}

	if result.Next() {
		user.Guid = result.Record().Values[0].(string)
		user.FIO = result.Record().Values[1].(string)
		user.Email = email
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
	result, err := session.Run("MATCH (n) WHERE n.guid=$guid_node MATCH (u:User) WHERE NOT (u)-[:access]->(n) AND NOT (:Admin)-[:role]->(u) RETURN u", map[string]interface{}{
		"guid_node": guid_node,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(guid_node)
	for result.Next() {
		props := result.Record().Values[0].(neo4j.Node).Props

		mas_user = append(mas_user, model.User{
			FIO:  props["FIO"].(string),
			Guid: props["guid"].(string),
		})
	}

	if err := result.Err(); err != nil {
		return nil, err
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

func (r *AuthRepository) IssueSessionUser(user model.User, refreshToken string) error {
	session := GetSession(*r.driver)
	defer session.Close()

	result, err := session.Run("MATCH (s:Session)-[:session]->(u:User) WHERE u.guid = $guid_user RETURN s", map[string]interface{}{
		"guid_user": user.Guid,
	})
	if err != nil {
		return err
	}

	if !result.Next() {
		_, err := session.Run("MATCH (u:User) WHERE u.guid = $guid_user CREATE (s:Session {refreshToken: $refreshToken}) CREATE (s)-[:session]->(u)", map[string]interface{}{
			"guid_user":    user.Guid,
			"refreshToken": refreshToken,
		})
		if err != nil {
			return err
		}
	} else {
		_, err := session.Run("MATCH (s)-[:session]->(u:User) WHERE u.guid = $guid_user SET s.refreshToken = $refreshToken", map[string]interface{}{
			"guid_user":    user.Guid,
			"refreshToken": refreshToken,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *AuthRepository) GetUserToRefreshToken(refreshToken string) (model.User, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	result, err := session.Run("MATCH (s:Session)-[:session]-(u:User) WHERE s.refreshToken = $refreshToken RETURN u.email, u.password", map[string]interface{}{
		"refreshToken": refreshToken,
	})
	if err != nil {
		return model.User{}, err
	}

	if result.Next() {
		return r.GetUser(result.Record().Values[0].(string), result.Record().Values[0].(string))
	} else {
		return model.User{}, ErrRecordNotFound
	}
}

func (r *AuthRepository) CreateResetPassword(user model.User) (string, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	guid := uuid.New().String()
	_, err := session.Run("CREATE (r:resetPassword {guid:$guid, code:$code})", map[string]interface{}{
		"guid": guid,
		"code": strings.Split(guid, "-")[0],
	})
	if err != nil {
		return "", err
	}
	_, err = session.Run("MATCH (u:User), (r:resetPassword) WHERE u.guid = $user_guid and r.guid = $guid CREATE (r)-[:reset]->(u)", map[string]interface{}{
		"guid":      guid,
		"user_guid": user.Guid,
	})
	if err != nil {
		return "", err
	}

	return strings.Split(guid, "-")[0], err
}

func (r *AuthRepository) CheckResetPassword(code string, user model.User) error {
	session := GetSession(*r.driver)
	defer session.Close()

	result, err := session.Run("MATCH (r:resetPassword)-[:reset]->(u:User) WHERE u.guid=$user_guid and r.code=$code RETURN r.guid", map[string]interface{}{
		"user_guid": user.Guid,
		"code":      code,
	})
	if err != nil {
		return err
	}

	if !result.Next() {
		return ErrRecordNotFound
	} else {
		return nil
	}
}

func (r *AuthRepository) UserResetPassword(user model.User, newPassword string) error {
	session := GetSession(*r.driver)
	defer session.Close()

	_, err := session.Run("MATCH (r:resetPassword)-[rr:reset]->(u:User) WHERE u.guid=$user_guid DELETE rr DELETE r", map[string]interface{}{
		"user_guid": user.Guid,
	})
	if err != nil {
		return err
	}
	_, err = session.Run("MATCH (u:User) WHERE u.guid=$user_guid SET u.Password = $newPassword", map[string]interface{}{
		"user_guid":   user.Guid,
		"newPassword": newPassword,
	})
	if err != nil {
		return err
	}

	return nil
}

package repository

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type RoleRepository struct {
	driver *neo4j.Driver
}

func NewRoleRepository(driver *neo4j.Driver) *RoleRepository {
	return &RoleRepository{driver: driver}
}

func (r *RoleRepository) IssueAccess(guid_user, guid_node string) (string, error) {
	session := GetSession(*r.driver)
	defer session.Close()
	fmt.Println(guid_user, guid_node)

	result, err := session.Run("MATCH (u:User),(n) WHERE u.guid = $guid_user AND n.guid = $guid_node CREATE (u)-[:access]->(n) RETURN labels(n)", map[string]interface{}{
		"guid_user": guid_user,
		"guid_node": guid_node,
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

func (r *RoleRepository) CheckRoleAdmin(guid_user string) (bool, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	result, err := session.Run("MATCH (:Admin)-[]->(u:User) WHERE u.guid = $guid_user RETURN u.guid", map[string]interface{}{
		"guid_user": guid_user,
	})
	if err != nil {
		return false, err
	}

	if result.Next() {
		return true, nil
	} else {
		return false, nil
	}
}

func (r *RoleRepository) CheckAccess(guid_user, guid_node string) (bool, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	result, err := session.Run("MATCH (u:User)-[y:access]->(n) WHERE u.guid = $guid_user AND n.guid=$guid_node RETURN y", map[string]interface{}{
		"guid_user": guid_user,
		"guid_node": guid_node,
	})
	if err != nil {
		return false, err
	}

	if result.Next() {
		return true, nil
	} else {
		return false, nil
	}
}

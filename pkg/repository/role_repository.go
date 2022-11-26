package repository

import (
	"github.com/google/uuid"
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

	result, err := session.Run("MATCH (u:User),(n) WHERE u.guid = $guid_user AND n.guid = $guid_node CREATE (u)-[:access]->(n) RETURN distinct labels(n)", map[string]interface{}{
		"guid_user": guid_user,
		"guid_node": guid_node,
	})
	if err != nil {
		return "", err
	}

	if result.Next() {
		for _, value := range result.Record().Values[0].([]interface{}) {
			return value.(string), nil
		}
	}
	return "", ErrRecordNotFound
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

func (r *RoleRepository) CreateInvite(guid_node string) (string, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	guid := uuid.New().String()
	result, err := session.Run("CREATE (i: Invite {guid:$guid, guidNode:$guidNode}) RETURN i.guid", map[string]interface{}{
		"guid":     guid,
		"guidNode": guid_node,
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

func (r *RoleRepository) UseInvite(guid_invite, guid_user string) error {
	session := GetSession(*r.driver)
	defer session.Close()

	guid_node, err := r.GetNodeToInvite(guid_invite)
	if err != nil {
		return err
	}
	_, err = r.IssueAccess(guid_user, guid_node)
	if err != nil {
		return err
	}

	_, err = session.Run("MATCH (i: Invite) WHERE i.guid = $guidInvite DELETE i", map[string]interface{}{
		"guidInvite": guid_invite,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *RoleRepository) GetNodeToInvite(guid_invite string) (string, error) {
	session := GetSession(*r.driver)
	defer session.Close()

	result, err := session.Run("MATCH (i:Invite) WHERE i.guid = $guidInvite RETURN i.guidNode", map[string]interface{}{
		"guidInvite": guid_invite,
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

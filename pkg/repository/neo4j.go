package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Config struct {
	URI      string
	Username string
	Password string
	DBName   string
}

func NewNeo4jDriver(cfg Config) (*neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(cfg.URI, neo4j.BasicAuth(cfg.Username, cfg.Password, ""))
	if err != nil {
		return nil, err
	}
	if err := driver.VerifyConnectivity(); err != nil {
		return nil, err
	}

	return &driver, nil
}

func GetSession(d neo4j.Driver) neo4j.Session {
	return d.NewSession(neo4j.SessionConfig{
		DatabaseName: "coursebook2",
	})
}

package postgres

import (
	"database/sql"
	"github.com/dzrry/activitycounter/vk/api"
	_ "github.com/lib/pq"
)

type Datastore interface {
	AllGroupMembers(int) (int, []*api.User, error)
	InsertGroupMembers(int, []*api.User) error
}

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

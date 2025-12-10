package storage

import (
	"database/sql"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

type Storage interface {
	CheckStorageExists() (bool, error)
	Init(password string) error
	Connect(password string) error
	CloseConnections() error
	GetConnection() *sql.DB
	PingCheck() bool
}

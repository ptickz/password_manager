package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"password_manager/internal/storage"
	"password_manager/internal/utils"
)

type SQLite struct {
	StorageUsername string
	StorageFilePath string
	MigrationPath   string
	Conn            *sql.DB
}

func NewSQLite(storageUsername, storageFilePath string, migrationPath string) storage.Storage {
	return &SQLite{
		StorageUsername: storageUsername,
		StorageFilePath: storageFilePath,
		MigrationPath:   migrationPath,
	}
}

func (s *SQLite) CheckStorageExists() (bool, error) {
	if _, err := os.Stat(s.StorageFilePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *SQLite) Init(password string) error {
	db, err := sql.Open(
		"sqlite3",
		s.buildConnString(password),
	)
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	if err != nil {
		return err
	}
	absMigrationPath := utils.GetAbsolutePath(s.MigrationPath)
	var sqlStatement []byte
	sqlStatement, err = os.ReadFile(absMigrationPath)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(sqlStatement))
	if err != nil {
		return err
	}
	return nil
}
func (s *SQLite) Connect(password string) error {
	db, err := sql.Open(
		"sqlite3",
		s.buildConnString(password),
	)
	if err != nil {
		log.Fatal(err)
	}
	s.Conn = db
	return nil
}

func (s *SQLite) CloseConnections() error {
	err := s.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLite) buildConnString(password string) string {
	connString := "file:" + s.StorageFilePath +
		"?_auth&_auth_user=" + s.StorageUsername +
		"&_auth_pass=" + password + "&_auth_crypt=SHA256"
	return connString
}

func (s *SQLite) GetConnection() *sql.DB {
	return s.Conn
}

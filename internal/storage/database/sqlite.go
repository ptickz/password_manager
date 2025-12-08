package database

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/mutecomm/go-sqlcipher/v4"
	"log"
	"net/url"
	"os"
	"password_manager/internal/storage"
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
	src := []byte(password)
	encodedKey := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(encodedKey, src)
	encodedKeyString := url.QueryEscape(string(encodedKey))
	db, err := sql.Open(
		"sqlite3",
		s.buildConnString(encodedKeyString),
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
	if _, err = db.Exec(fmt.Sprintf("PRAGMA key = x%s;", encodedKeyString)); err != nil {
		return err
	}

	var sqlStatement []byte
	sqlStatement, err = os.ReadFile(s.MigrationPath)
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

func (s *SQLite) buildConnString(encodedKeyString string) string {
	connString := s.StorageFilePath +
		fmt.Sprintf("?_pragma_key=x%s&_pragma_cipher_page_size=4096", encodedKeyString)
	return connString
}

func (s *SQLite) GetConnection() *sql.DB {
	return s.Conn
}

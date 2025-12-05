package config

import (
	"github.com/joho/godotenv"
	"os"
)

const envFilePath = "/config/.env"

type Config struct {
	Timeout         string
	StorageType     string
	StorageFilePath string
	StorageUsername string
	MigrationPath   string
	TransportType   string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	c := Config{
		Timeout:         os.Getenv("TIMEOUT"),
		StorageFilePath: os.Getenv("STORAGE_FILE_PATH"),
		StorageUsername: os.Getenv("STORAGE_USERNAME"),
		MigrationPath:   os.Getenv("MIGRATION_PATH"),
		StorageType:     os.Getenv("STORAGE_TYPE"),
		TransportType:   os.Getenv("TRANSPORT_TYPE"),
	}

	return &c, nil
}

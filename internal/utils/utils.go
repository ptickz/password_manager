package utils

import (
	"os"
	"path/filepath"
)

func GetAbsolutePath(relativePath string) string {
	return filepath.Join(os.Getenv("APP_PATH"), os.Getenv("APP_NAME"), relativePath)
}

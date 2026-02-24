package utils

import (
	"os"
	"path/filepath"
)

func EnsureFileExists(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return os.WriteFile(filePath, []byte("[]"), 0644)
	}
	return nil
}

package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadJSON[T any](filePath string) ([]T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []T{}, nil
		}
		return nil, fmt.Errorf("could not open file %s: %w", filePath, err)
	}
	defer file.Close()
	var data []T
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil && err.Error() != "EOF" {
		return nil, fmt.Errorf("failed to decode JSON from %s: %w", filePath, err)
	}
	if data == nil {
		return []T{}, nil
	}
	return data, nil
}

func WriteJSON[T any](filePath string, data []T) error {
	fileData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	return os.WriteFile(filePath, fileData, 0644)
}

func ExtractID(path string, prefix string) string {
	if len(path) <= len(prefix) {
		return ""
	}
	return path[len(prefix):]
}

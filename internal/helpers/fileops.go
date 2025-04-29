package helpers

import (
	"fmt"
	"os"
)

func ReadFileAsString(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		LogError(fmt.Errorf("ReadFile(): Error while reading file: %w", err))
	}
	return string(content)
}

func List(path string) []os.DirEntry {
	entries, err := os.ReadDir(path)
	if err != nil {
		LogError(fmt.Errorf("ListFiles(): Error while reading directory: %w", err))
	}
	return entries
}

func ListAsString(path string) []string {
	dirEntries := List(path)
	var items []string

	for _, entry := range dirEntries {
		items = append(items, entry.Name())
	}
	return items
}

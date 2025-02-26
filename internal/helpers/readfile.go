package helpers

import (
	"log"
	"os"
)

func ReadFile(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

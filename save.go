package main

import (
	"log"
	"os"

	"github.com/owulveryck/gptslideshow/config"
)

func saveContent(filename string, content []byte) error {
	// Create a temporary file within the directory
	tempFile, err := os.CreateTemp(config.ConfigInstance.TempDir, filename)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// Write the string to the file
	_, err = tempFile.Write(content)
	if err != nil {
		return err
	}

	// Output the temporary file path
	log.Printf("Temporary file created: %s", tempFile.Name())

	return nil
}

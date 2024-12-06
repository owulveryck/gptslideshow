package main

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/api/drive/v3"
)

func main() {
	client := initGoogleClient()

	// Handle template copy if specified
	srvDrive := initDriveService(client)

	// Open the image file to upload
	file, err := os.Open("testdata/sample.png")
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer file.Close()

	// Create a drive.File instance with the desired metadata
	driveFile := &drive.File{
		Name: "testdata/sample.png",
	}

	// Upload the file to Google Drive
	uploadedFile, err := srvDrive.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		log.Fatalf("Unable to upload file: %v", err)
	}

	fmt.Printf("File uploaded successfully. File ID: %s\n", uploadedFile.Id)

	// Make the file public by setting the permission
	_, err = srvDrive.Permissions.Create(uploadedFile.Id, &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}).Do()
	if err != nil {
		log.Fatalf("Unable to set file permission: %v", err)
	}

	// Construct the public URL
	publicURL := fmt.Sprintf("https://drive.google.com/uc?id=%s", uploadedFile.Id)
	fmt.Printf("Public URL: %s\n", publicURL)
}

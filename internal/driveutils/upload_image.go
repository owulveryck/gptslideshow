package driveutils

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"log"

	"google.golang.org/api/drive/v3"
)

// Upload an image on google drive and returns the public url
func UploadImage(ctx context.Context, srvDrive *drive.Service, img image.Image, name string) (string, error) {
	// Create a drive.File instance with the desired metadata
	driveFile := &drive.File{
		Name: name,
	}
	// Encodez l'image en PNG dans un buffer
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatalf("failed to encode image: %v", err)
	}

	// Upload the file to Google Drive
	uploadedFile, err := srvDrive.Files.Create(driveFile).Media(&buf).Do()
	if err != nil {
		return "", err
	}

	fmt.Printf("File uploaded successfully. File ID: %s\n", uploadedFile.Id)

	// Make the file public by setting the permission
	_, err = srvDrive.Permissions.Create(uploadedFile.Id, &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}).Do()
	if err != nil {
		return "", err
	}

	// Construct the public URL
	publicURL := fmt.Sprintf("https://drive.google.com/uc?id=%s", uploadedFile.Id)
	return publicURL, nil
}

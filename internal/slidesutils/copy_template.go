package slidesutils

import (
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
)

// CopyTemplate copies a presentation template and returns the new presentation ID.
func CopyTemplate(ctx context.Context, driveSrv *drive.Service, templatePresentationId string) (string, error) {
	copyTitle := "gptSlides"
	copyRequest := &drive.File{Name: copyTitle}
	copiedFile, err := driveSrv.Files.Copy(templatePresentationId, copyRequest).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("unable to copy presentation: %v", err)
	}
	fmt.Printf("Copied presentation ID: %s\n", copiedFile.Id)
	return copiedFile.Id, nil
}

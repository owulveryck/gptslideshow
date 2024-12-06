package main

import (
	"context"
	"log"

	"github.com/owulveryck/gptslideshow/internal/slidesutils"
	drive "google.golang.org/api/drive/v3"
)

func handleTemplateCopy(ctx context.Context, driveSrv *drive.Service, templateId string) string {
	presentationId, err := slidesutils.CopyTemplate(ctx, driveSrv, templateId)
	if err != nil {
		log.Fatalf("Unable to copy presentation: %v", err)
	}
	return presentationId
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	slides "google.golang.org/api/slides/v1"

	"github.com/owulveryck/gptslideshow/internal/gcputils"
)

func initGoogleClient() *http.Client {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope, slides.PresentationsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return gcputils.GetClient(config)
}

func initSlidesService(ctx context.Context, client *http.Client) *slides.Service {
	// slidesSrv, err := slides.New(client)
	slidesSrv, err := slides.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}
	return slidesSrv
}

func initDriveService(ctx context.Context, client *http.Client) *drive.Service {
	driveSrv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return driveSrv
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/owulveryck/gptslideshow/internal/gcputils"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	slides "google.golang.org/api/slides/v1"
)

// customRoundTripper adds a custom header to all requests.
type customRoundTripper struct {
	original http.RoundTripper
}

// RoundTrip implements the RoundTripper interface.
func (c *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original.
	clonedReq := req.Clone(req.Context())
	// Add the ngrok header
	clonedReq.Header.Set("ngrok-skip-browser-warning", "true")

	// Use the original RoundTripper to send the request.
	return c.original.RoundTrip(clonedReq)
}

func initGoogleClientPerso(credentials string) *http.Client {
	b, err := os.ReadFile(credentials)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope, slides.PresentationsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := gcputils.GetClientPerso(config)
	client.Transport = &customRoundTripper{
		original: client.Transport,
	}
	return client
}

func initGoogleClient(credentials string) *http.Client {
	b, err := os.ReadFile(credentials)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope, slides.PresentationsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := gcputils.GetClient(config)
	client.Transport = &customRoundTripper{
		original: client.Transport,
	}
	return client
}

func initSlidesService(ctx context.Context, client *http.Client) *slides.Service {
	// slidesSrv, err := slides.New(client)
	slidesSrv, err := slides.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}
	return slidesSrv
}

func initDriveService(client *http.Client) *drive.Service {
	driveSrv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return driveSrv
}

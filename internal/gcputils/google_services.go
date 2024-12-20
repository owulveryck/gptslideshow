// Package gcputils provides utility functions for initializing Google API clients
// and services such as Google Drive and Google Slides. It simplifies the
// authentication process and service initialization using credentials and token caching.
package gcputils

import (
	"context"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	slides "google.golang.org/api/slides/v1"
)

// InitGoogleClient initializes an HTTP client for accessing Google APIs using
// credentials from a JSON file and a cached token file. If the token is missing
// or expired, the function triggers the authentication process to obtain a new
// token, which is then saved for future use.
//
// Parameters:
//   - credentials: Path to the JSON file containing the client credentials.
//   - tokenCacheFile: Path to the file where the authentication token is cached.
//
// Returns:
//
//	*http.Client: An HTTP client configured for use with Google APIs.
//
// Errors:
//   - Logs a fatal error if the credentials file cannot be read or parsed.
func InitGoogleClient(credentials, tokenCacheFile string) *http.Client {
	b, err := os.ReadFile(credentials)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope, slides.PresentationsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return getClient(config, tokenCacheFile)
}

// InitSlidesService initializes a Google Slides service using the provided
// context and HTTP client.
//
// Parameters:
//   - ctx: The context for the service.
//   - client: The HTTP client for authenticating requests.
//
// Returns:
//
//	*slides.Service: A service for interacting with Google Slides.
//
// Errors:
//   - Logs a fatal error if the Slides service cannot be created.
func InitSlidesService(ctx context.Context, client *http.Client) *slides.Service {
	slidesSrv, err := slides.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}
	return slidesSrv
}

// InitDriveService initializes a Google Drive service using the provided
// context and HTTP client.
//
// Parameters:
//   - ctx: The context for the service.
//   - client: The HTTP client for authenticating requests.
//
// Returns:
//
//	*drive.Service: A service for interacting with Google Drive.
//
// Errors:
//   - Logs a fatal error if the Drive service cannot be created.
func InitDriveService(ctx context.Context, client *http.Client) *drive.Service {
	driveSrv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return driveSrv
}

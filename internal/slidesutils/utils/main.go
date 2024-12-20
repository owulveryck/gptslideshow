package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"

	"github.com/owulveryck/gptslideshow/internal/gcputils"
)

func main() {
	var presentationId string
	flag.StringVar(&presentationId, "id", "", "ID of the slide to update, empty means create a new one")

	flag.Parse()

	ctx := context.Background()
	// Initialize Google services
	cacheDir, err := gcputils.GetCredentialCacheDir()
	if err != nil {
		log.Fatal(err)
	}

	client := gcputils.InitGoogleClient("../../../credentials.json", filepath.Join(cacheDir, "slides.googleapis.com-go-quickstart.json"))
	slidesSrv := gcputils.InitSlidesService(ctx, client)

	err = ListSlidesAndElements(slidesSrv, presentationId)
	if err != nil {
		log.Fatal(err)
	}
}

// Package ngrok provides functionality to expose a locally hosted image
// to the public internet using ngrok tunneling.
//
// This package implements the uploader.Uploader interface, allowing seamless
// integration with systems that require a publicly accessible URL for images.
package ngrok

import (
	"context"
	"image"
	"image/png"
	"log"
	"net/http"
	"path"
	"strings"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

// Uploader is a struct that manages image storage and provides public URLs
// for images through an ngrok tunnel.
//
// It implements the uploader.Uploader interface, enabling compatibility with
// systems that require image upload and public access capabilities.
type Uploader struct {
	imageStore map[string]image.Image // In-memory storage for uploaded images
	tunURL     string                 // Public URL of the ngrok tunnel
}

// NewUploader creates and initializes a new Uploader instance.
//
// The Uploader uses an in-memory map to store images and can be configured
// to expose these images through a public URL using an ngrok tunnel.
//
// Returns:
//   - A pointer to the newly created Uploader instance.
func NewUploader() *Uploader {
	return &Uploader{
		imageStore: make(map[string]image.Image),
	}
}

// Upload uploads the provided image with the specified name to a destination
// and returns a publicly accessible URL for the uploaded image.
//
// Parameters:
//
//	ctx: The context for managing request deadlines, cancellation signals,
//	     and other request-scoped values.
//	name: The desired name or identifier for the uploaded image.
//	img: The image.Image object representing the image to upload.
//
// Returns:
//
//	string: A publicly accessible URL where the uploaded image can be accessed.
//	error: An error object if the upload fails, or nil if the operation is successful.
func (uploader *Uploader) Upload(ctx context.Context, name string, img image.Image) (string, error) {
	uploader.imageStore[name] = img
	return path.Join(uploader.tunURL, name), nil
}

// StartServer starts an ngrok tunnel to expose the local HTTP server to the
// public internet. The public URL is stored in the Uploader for later use.
// This should be launched within a goroutine as this function does not return
//
// Parameters:
//   - ctx: The context to manage the lifecycle of the ngrok tunnel.
//
// Returns:
//   - An error if the tunnel fails to start.
func (uploader *Uploader) StartServer(ctx context.Context) {
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		log.Fatal("Error starting Ngrok tunnel:", err)
	}
	uploader.tunURL = listener.URL()

	// Set up the HTTP server
	http.Serve(listener, http.HandlerFunc(uploader.imageHandler))
}

func (uploader *Uploader) imageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the image name from the URL path
	imageName := strings.TrimLeft(r.URL.Path, "/")
	if imageName == "" {
		http.Error(w, "Image name is required", http.StatusBadRequest)
		return
	}

	// Find the image in the store
	img, exists := uploader.imageStore[imageName]
	if !exists {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// Set the content type to image/png
	w.Header().Set("Content-Type", "image/png")

	// Encode the image as PNG and write to the response
	err := png.Encode(w, img)
	if err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}
}

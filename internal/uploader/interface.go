// Package uploader provides an abstraction for uploading images to a destination
// and retrieving a publicly accessible URL for the uploaded content.
package uploader

import (
	"context"
	"image"
)

// Uploader defines an interface for uploading images to a destination.
// Implementations of this interface should ensure the uploaded images
// are accessible via a publicly available URL.
type Uploader interface {
	// Upload uploads the provided image with the specified name to a destination
	// and returns a publicly accessible URL for the uploaded image.
	//
	// Parameters:
	//   ctx: The context for managing request deadlines, cancellation signals,
	//        and other request-scoped values.
	//   name: The desired name or identifier for the uploaded image.
	//   img: The image.Image object representing the image to upload.
	//
	// Returns:
	//   string: A publicly accessible URL where the uploaded image can be accessed.
	//   error: An error object if the upload fails, or nil if the operation is successful.
	Upload(ctx context.Context, name string, img image.Image) (string, error)
}

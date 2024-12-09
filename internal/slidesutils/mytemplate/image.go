package mytemplate

import (
	"context"
	"fmt"

	slides "google.golang.org/api/slides/v1"
)

// InsertImage inserts an image into the current slide of the presentation.
// The image is positioned and sized according to specified parameters.
//
// Parameters:
//   - ctx: A context to manage request lifetime.
//   - imageUrl: The URL of the image to be inserted.
//   - width: The width of the image in EMUs.
//   - height: The height of the image in EMUs.
//   - translateX: The X translation of the image in EMUs.
//   - translateY: The Y translation of the image in EMUs.
//
// Returns:
//   - error: An error if the image insertion fails.
func (b *Builder) InsertImage(ctx context.Context, imageUrl string, width, height, translateX, translateY float64) error {
	if b.CurrentSlide == nil {
		return fmt.Errorf("current slide is not set")
	}

	// Define the properties for the image
	imageRequest := &slides.Request{
		CreateImage: &slides.CreateImageRequest{
			Url: imageUrl,
			ElementProperties: &slides.PageElementProperties{
				PageObjectId: b.CurrentSlide.ObjectId,
				Size: &slides.Size{
					Height: &slides.Dimension{
						Magnitude: height,
						Unit:      "EMU",
					},
					Width: &slides.Dimension{
						Magnitude: width,
						Unit:      "EMU",
					},
				},
				Transform: &slides.AffineTransform{
					ScaleX:     1.0,
					ScaleY:     1.0,
					TranslateX: translateX,
					TranslateY: translateY,
					Unit:       "EMU",
				},
			},
		},
	}

	// Create a batch update request
	batchUpdateRequest := &slides.BatchUpdatePresentationRequest{
		Requests: []*slides.Request{imageRequest},
	}

	// Execute the batch update request
	_, err := b.Srv.Presentations.BatchUpdate(b.Presentation.PresentationId, batchUpdateRequest).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to insert image: %w", err)
	}

	return nil
}

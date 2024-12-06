/*
Package slidesutils provides utilities for working with Google Slides presentations,
including creating and managing slides programmatically using the Google Slides API.

This package is intended to streamline common tasks such as creating slides, managing chapters,
and interacting with Google Slides presentation objects.
*/
package slidesutils

import (
	"context"

	slides "google.golang.org/api/slides/v1"
)

// Builder encapsulates state and methods for working with a Google Slides presentation.
// It maintains the current chapter, the current slide, and the overall presentation object.
type Builder struct {
	Srv            *slides.Service      // The Google Slides API service client.
	CurrentChapter int                  // Tracks the current chapter number in the presentation.
	CurrentSlide   *slides.Page         // Points to the current slide being manipulated.
	Presentation   *slides.Presentation // The full presentation being managed.
	slideNumber    int
}

const (
	// ChapterLayoutId represents the layout ID for a chapter slide in the presentation.
	ChapterLayoutId = "g2ac55f3490c_0_1010"

	// TitleSubtitleBody represents the layout ID for a slide with a title, subtitle, and body content.
	TitleSubtitleBody = "g2ac55f3490c_0_1006"
)

// NewBuilder initializes a Builder instance for managing a Google Slides presentation.
//
// This function retrieves the presentation with the specified ID using the provided Google Slides API service client.
// The resulting Builder is used for creating and managing slides programmatically.
//
// Parameters:
//   - ctx: A context to manage request lifetime.
//   - srv: A Google Slides API service client.
//   - presentationId: The ID of the Google Slides presentation to manage.
//
// Returns:
//   - *Builder: A new Builder instance for the specified presentation filled with the Srv and Presentation field.
//   - error: An error if the presentation could not be retrieved or if the API call fails.
func NewBuilder(ctx context.Context, srv *slides.Service, presentationId string) (*Builder, error) {
	presentation, err := srv.Presentations.Get(presentationId).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return &Builder{
		Srv:            srv,
		CurrentChapter: 0,
		CurrentSlide:   nil,
		Presentation:   presentation,
	}, nil
}

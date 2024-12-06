package slidesutils

import (
	"context"
	"fmt"

	"github.com/owulveryck/gptslideshow/internal/structure"
	slides "google.golang.org/api/slides/v1"
)

// CreateSlideTitleSubtitleBody creates a new slide with a title, subtitle, and body content.
// It uses the predefined layout for title, subtitle, and body.
//
// Parameters:
//   - ctx: A context to manage request lifetime.
//   - slide: A structure containing slide information such as title, subtitle, and body.
//
// Returns:
//   - error: An error if the slide creation or text insertion fails.
func (b *Builder) CreateSlideTitleSubtitleBody(ctx context.Context, slide structure.Slide) error {
	// Use the CreateNewSlide method to create a new slide with the title, subtitle, and body layout.
	if err := b.CreateNewSlide(ctx, TitleSubtitleBody); err != nil {
		return fmt.Errorf("failed to create slide with title, subtitle, and body: %w", err)
	}

	// Ensure the current slide is set after creation.
	if b.CurrentSlide == nil {
		return fmt.Errorf("current slide is not set after creation")
	}

	// Find placeholders for title, subtitle, and body in the newly created slide.
	var titlePlaceholderID, subtitlePlaceholderID, bodyPlaceholderID string
	for _, element := range b.CurrentSlide.PageElements {
		if element.Shape != nil && element.Shape.Placeholder != nil {
			switch element.Shape.Placeholder.Type {
			case "TITLE":
				titlePlaceholderID = element.ObjectId
			case "SUBTITLE":
				subtitlePlaceholderID = element.ObjectId
			case "BODY":
				bodyPlaceholderID = element.ObjectId
			}
		}
	}

	// Check if all placeholders were found.
	if titlePlaceholderID == "" || bodyPlaceholderID == "" || subtitlePlaceholderID == "" {
		return fmt.Errorf("failed to find placeholders on the new slide")
	}

	// Prepare text requests to insert the title, subtitle, and body content.
	textRequests := []*slides.Request{
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       titlePlaceholderID,
				InsertionIndex: 0,
				Text:           slide.Title,
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       subtitlePlaceholderID,
				InsertionIndex: 0,
				Text:           slide.Subtitle,
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       bodyPlaceholderID,
				InsertionIndex: 0,
				Text:           slide.Body,
			},
		},
	}

	// Execute the batch update request to insert text into the placeholders.
	if _, err := b.Srv.Presentations.BatchUpdate(b.Presentation.PresentationId, &slides.BatchUpdatePresentationRequest{
		Requests: textRequests,
	}).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to insert text: %w", err)
	}

	return nil
}

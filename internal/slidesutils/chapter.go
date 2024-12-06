package slidesutils

import (
	"context"
	"fmt"
	"strconv"

	"github.com/owulveryck/gptslideshow/internal/structure"
	slides "google.golang.org/api/slides/v1"
)

// CreateChapter creates a new chapter slide in the presentation.
// It uses the predefined chapter layout and updates the slide with the chapter title and number.
//
// Parameters:
//   - ctx: A context to manage request lifetime.
//   - slide: A structure containing slide information such as title.
//
// Returns:
//   - error: An error if the slide creation or text insertion fails.
func (b *Builder) CreateChapter(ctx context.Context, slide structure.Slide) error {
	// Use the CreateNewSlide method to create a new slide with the chapter layout.
	if err := b.CreateNewSlide(ctx, ChapterLayoutId); err != nil {
		return fmt.Errorf("failed to create chapter slide: %w", err)
	}

	// Ensure the current slide is set after creation.
	if b.CurrentSlide == nil {
		return fmt.Errorf("current slide is not set after creation")
	}

	// Find placeholders for title and body in the newly created slide.
	var titlePlaceholderID, bodyPlaceholderID string
	for _, element := range b.CurrentSlide.PageElements {
		if element.Shape != nil && element.Shape.Placeholder != nil {
			switch element.Shape.Placeholder.Type {
			case "TITLE":
				titlePlaceholderID = element.ObjectId
			case "BODY":
				bodyPlaceholderID = element.ObjectId
			}
		}
	}

	// Check if both placeholders were found.
	if titlePlaceholderID == "" || bodyPlaceholderID == "" {
		return fmt.Errorf("failed to find placeholders on the new slide")
	}

	// Prepare text requests to insert the chapter title and number.
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
				ObjectId:       bodyPlaceholderID,
				InsertionIndex: 0,
				Text:           strconv.Itoa(b.CurrentChapter),
			},
		},
	}

	// Execute the batch update request to insert text into the placeholders.
	if _, err := b.Srv.Presentations.BatchUpdate(b.Presentation.PresentationId, &slides.BatchUpdatePresentationRequest{
		Requests: textRequests,
	}).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to insert text: %w", err)
	}

	// Increment the current chapter number after successful slide creation.
	b.CurrentChapter++
	return nil
}

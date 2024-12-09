package mytemplate

import (
	"context"
	"fmt"
	"log"
	"time"

	slides "google.golang.org/api/slides/v1"
)

// CreateCover creates a new cover slide in the presentation.
// It uses the predefined chapter layout and updates the slide with the chapter title and number.
//
// Parameters:
//   - ctx: A context to manage request lifetime.
//   - slide: A structure containing slide information such as title.
//
// Returns:
//   - error: An error if the slide creation or text insertion fails.
func (b *Builder) CreateCover(ctx context.Context, title, subtitle string) error {
	// Use the CreateNewSlide method to create a new slide with the chapter layout.
	if err := b.CreateNewSlide(ctx, CoverLayoutID); err != nil {
		return fmt.Errorf("failed to create cover slide: %w", err)
	}

	// Ensure the current slide is set after creation.
	if b.CurrentSlide == nil {
		return fmt.Errorf("current slide is not set after creation")
	}

	// Find placeholders for title and body in the newly created slide.
	titlesID := make([]string, 0, 3)
	var bodyPlaceholderID string
	for _, element := range b.CurrentSlide.PageElements {
		log.Println(element.Shape.Placeholder.Type)
		if element.Shape != nil && element.Shape.Placeholder != nil {
			switch element.Shape.Placeholder.Type {
			case "TITLE":
				titlesID = append(titlesID, element.ObjectId)
			case "SUBTITLE":
				bodyPlaceholderID = element.ObjectId
			}
		}
	}
	// Get the current date
	currentDate := time.Now()
	// Prepare text requests to insert the chapter title and number.
	textRequests := []*slides.Request{
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       titlesID[0],
				InsertionIndex: 0,
				Text:           title,
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       titlesID[1],
				InsertionIndex: 0,
				Text:           currentDate.Format("01/02/2006"),
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       titlesID[2],
				InsertionIndex: 0,
				Text:           "gptSlideShow",
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       bodyPlaceholderID,
				InsertionIndex: 0,
				Text:           subtitle,
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

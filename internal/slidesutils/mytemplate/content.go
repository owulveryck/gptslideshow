package mytemplate

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/owulveryck/gptslideshow/internal/slidesutils"
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

	formattedBody := slidesutils.InsertMarkdownContent(slide.Body, bodyPlaceholderID)
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
	}

	textRequests = append(textRequests, formattedBody...)
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(textRequests)

	// Sort the requests slice
	sort.SliceStable(textRequests, func(i, j int) bool {
		// Define the priority for each type
		priority := func(req *slides.Request) int {
			switch {
			case req.InsertText != nil:
				return 1 // Highest priority
			case req.UpdateTextStyle != nil:
				return 2
			case req.CreateParagraphBullets != nil:
				return 3 // Lowest priority
			default:
				return 4 // Fallback for unknown types
			}
		}

		// Compare the priorities of the two elements
		return priority(textRequests[i]) < priority(textRequests[j])
	})
	enc.Encode(textRequests)

	for _, textRequests := range textRequests {
		// Execute the batch update request to insert text into the placeholders.
		if _, err := b.Srv.Presentations.BatchUpdate(b.Presentation.PresentationId, &slides.BatchUpdatePresentationRequest{
			Requests: []*slides.Request{textRequests},
		}).Context(ctx).Do(); err != nil {
			return fmt.Errorf("failed to insert text: %w", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

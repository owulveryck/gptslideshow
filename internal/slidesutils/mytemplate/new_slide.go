package mytemplate

import (
	"context"
	"fmt"

	slides "google.golang.org/api/slides/v1"
)

// CreateNewSlide creates a new slide in the presentation using the specified layout ID.
// It updates the Builder's CurrentSlide to reference the newly created slide.
//
// Parameters:
//   - ctx: A context to manage request lifetime.
//   - layoutId: The ID of the layout to be used for the new slide.
//
// Returns:
//   - error: An error if the slide creation fails or if the new slide cannot be matched
//     to the presentation's slide list.
func (b *Builder) CreateNewSlide(ctx context.Context, layoutId string) error {
	// Construct the request for creating a new slide using the specified layout ID.
	createSlideRequests := []*slides.Request{
		{
			CreateSlide: &slides.CreateSlideRequest{
				SlideLayoutReference: &slides.LayoutReference{
					LayoutId: layoutId,
				},
			},
		},
	}

	// Execute the batch update request to create a new slide.
	createSlideResponse, err := b.Srv.Presentations.BatchUpdate(b.Presentation.PresentationId, &slides.BatchUpdatePresentationRequest{
		Requests: createSlideRequests,
	}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create slide: %w", err)
	}

	// Extract the new slide ID from the response.
	var newSlideID string
	for _, reply := range createSlideResponse.Replies {
		if reply.CreateSlide != nil {
			newSlideID = reply.CreateSlide.ObjectId
			break
		}
	}

	// Check if a new slide ID was retrieved; return an error if not.
	if newSlideID == "" {
		return fmt.Errorf("failed to retrieve new slide ID")
	}

	// Refresh the presentation to include the newly created slide.
	presentation, err := b.Srv.Presentations.Get(b.Presentation.PresentationId).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to retrieve updated presentation: %w", err)
	}

	// Update the current slide reference in the Builder to the newly created slide.
	var found bool
	for _, slide := range presentation.Slides {
		if slide.ObjectId == newSlideID {
			b.CurrentSlide = slide
			found = true
			break
		}
	}

	// Ensure the new slide was matched to an existing slide in the presentation.
	if !found {
		return fmt.Errorf("new slide ID %q not found in the updated presentation", newSlideID)
	}

	// Update the Builder's Presentation to reflect the latest state.
	b.Presentation = presentation

	return nil
}

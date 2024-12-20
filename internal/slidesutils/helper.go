package slidesutils

import (
	"context"
	"iter"
	"log"
	"strings"

	"github.com/owulveryck/gptslideshow/internal/ai"
	"google.golang.org/api/slides/v1"
)

// Presentation allows interaction with a Google Slides presentation.
// It provides methods for updating the presentation, checking for changes,
// and iterating through the slides in the presentation.
type Presentation struct {
	PresentationID          string                            // The unique identifier of the presentation.
	PresentationService     *slides.PresentationsService      // Service for interacting with the Google Slides API.
	PresentationPageService *slides.PresentationsPagesService // Service for interacting with individual slides.
	presentation            *slides.Presentation              // Cached representation of the presentation.
	presentationRevision    string                            // The current revision ID of the presentation.
}

// UpdatePresentationPointer retrieves the latest version of the presentation
// from the Google Slides API and updates the cached presentation object.
func (p *Presentation) UpdatePresentationPointer() error {
	presentation, err := p.PresentationService.Get(p.PresentationID).Do()
	if err != nil {
		return err
	}
	p.presentation = presentation
	return nil
}

// Update applies a batch of update requests to the presentation.
// It takes a context, and performs the updates defined in the
// provided slice of Request objects.
func (p *Presentation) Update(ctx context.Context, requests []*slides.Request) error {
	// Execute the batch update
	_, err := p.PresentationService.BatchUpdate(p.PresentationID, &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}
	return nil
}

// HasChanged checks if the presentation has been updated since the last check.
// It compares the current revision ID of the presentation with the stored revision ID.
// If the revision has changed, it returns true and updates the stored revision ID.
func (p *Presentation) HasChanged() bool {
	if p.presentation.RevisionId == p.presentationRevision {
		return false
	}
	p.presentationRevision = p.presentation.RevisionId
	return true
}

// Slides returns an iterator for the slides in the current presentation.
// It iterates through the slides and yields each slide as a *slides.Page object.
// If an error occurs while retrieving a slide, it logs the error and continues iterating.
func (p *Presentation) Slides() iter.Seq[*slides.Page] {
	return func(yield func(*slides.Page) bool) {
		for _, s := range p.presentation.Slides {
			slide, err := p.PresentationPageService.Get(p.PresentationID, s.ObjectId).Do()
			if err != nil {
				log.Println(err)
			}
			if !yield(slide) {
				return
			}
		}
	}
}

func (p *Presentation) UpdateFromAI(ctx context.Context, ai ai.AIInterface, objectID, prompt, tagWord string) error {
	prompt = strings.ReplaceAll(prompt, tagWord, "")
	result, err := ai.SimpleQuery(ctx, prompt)
	if err != nil {
		return err
	}
	requests := processText(objectID, result)
	// Execute the batch update
	_, err = p.PresentationService.BatchUpdate(p.PresentationID, &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}
	return nil
}

func processText(objectID, input string) []*slides.Request {
	content := []*slides.Request{
		{
			DeleteText: &slides.DeleteTextRequest{
				ObjectId: objectID,
				TextRange: &slides.Range{
					Type: "ALL",
				},
			},
		},
	}

	content = append(content, InsertMarkdownContent(input, objectID)...)
	SortRequests(content)
	// Create a batch update request to replace the text
	return content
}

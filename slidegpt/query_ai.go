package main

import (
	"context"
	"iter"
	"log"
	"strings"

	"github.com/owulveryck/gptslideshow/internal/ai"
	slides "google.golang.org/api/slides/v1"
)

type helper struct {
	presentationID          string
	presentationService     *slides.PresentationsService
	presentation            *slides.Presentation
	presentationRevision    string
	presentationPageService *slides.PresentationsPagesService
}

func (h *helper) updatePresentationPointer() error {
	presentation, err := h.presentationService.Get(h.presentationID).Do()
	if err != nil {
		return err
	}
	h.presentation = presentation
	return nil
}

func (h *helper) updateFromAI(ctx context.Context, ai ai.AIInterface, objectID, prompt, tagWord string) error {
	prompt = strings.ReplaceAll(prompt, tagWord, "")
	result, err := ai.SimpleQuery(ctx, prompt)
	if err != nil {
		return err
	}
	requests := processText(objectID, result)
	// Execute the batch update
	_, err = h.presentationService.BatchUpdate(h.presentationID, &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}
	return nil
}

func (h *helper) presentationHasChanged() bool {
	if h.presentation.RevisionId == h.presentationRevision {
		return false
	}
	h.presentationRevision = h.presentation.RevisionId
	return true
}

// Create an iterator for slides in the current presentation
func (h *helper) Slides() iter.Seq[*slides.Page] {
	return func(yield func(*slides.Page) bool) {
		for _, s := range h.presentation.Slides {
			slide, err := h.presentationPageService.Get(h.presentationID, s.ObjectId).Do()
			if err != nil {
				log.Println(err)
			}
			if !yield(slide) {
				return
			}
		}
	}
}

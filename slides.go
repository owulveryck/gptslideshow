package main

import (
	"context"
	"fmt"

	drive "google.golang.org/api/drive/v3"
	slides "google.golang.org/api/slides/v1"
)

// CopyTemplate copies a presentation template and returns the new presentation ID.
func CopyTemplate(ctx context.Context, driveSrv *drive.Service, templatePresentationId string) (string, error) {
	copyTitle := "gptSlides"
	copyRequest := &drive.File{Name: copyTitle}
	copiedFile, err := driveSrv.Files.Copy(templatePresentationId, copyRequest).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("unable to copy presentation: %v", err)
	}
	fmt.Printf("Copied presentation ID: %s\n", copiedFile.Id)
	return copiedFile.Id, nil
}

// CreateSlide creates a new slide in the presentation.
func CreateSlide(ctx context.Context, srv *slides.Service, presentationId string, slide Slide) error {
	createSlideRequests := []*slides.Request{
		{
			CreateSlide: &slides.CreateSlideRequest{
				SlideLayoutReference: &slides.LayoutReference{
					LayoutId: "g2ac55f3490c_0_1006", // Replace with your custom layout ID
				},
			},
		},
	}

	createSlideResponse, err := srv.Presentations.BatchUpdate(presentationId, &slides.BatchUpdatePresentationRequest{
		Requests: createSlideRequests,
	}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create slide: %v", err)
	}

	var newSlideID string
	for _, reply := range createSlideResponse.Replies {
		if reply.CreateSlide != nil {
			newSlideID = reply.CreateSlide.ObjectId
			break
		}
	}

	if newSlideID == "" {
		return fmt.Errorf("failed to retrieve new slide ID")
	}

	presentation, err := srv.Presentations.Get(presentationId).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to retrieve presentation: %v", err)
	}

	var titlePlaceholderID, subtitlePlaceholderID, bodyPlaceholderID string
	for _, slide := range presentation.Slides {
		if slide.ObjectId == newSlideID {
			for _, element := range slide.PageElements {
				if element.Shape != nil && element.Shape.Placeholder != nil {
					switch element.Shape.Placeholder.Type {
					case "TITLE":
						titlePlaceholderID = element.ObjectId
					case "BODY":
						bodyPlaceholderID = element.ObjectId
					case "SUBTITLE":
						subtitlePlaceholderID = element.ObjectId
					}
				}
			}
			break
		}
	}

	if titlePlaceholderID == "" || bodyPlaceholderID == "" || subtitlePlaceholderID == "" {
		return fmt.Errorf("failed to find placeholders on the new slide")
	}

	textRequests := []*slides.Request{
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       subtitlePlaceholderID,
				InsertionIndex: 0,
				Text:           slide.Subtitle,
			},
		},
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
				Text:           slide.Body,
			},
		},
	}

	_, err = srv.Presentations.BatchUpdate(presentationId, &slides.BatchUpdatePresentationRequest{
		Requests: textRequests,
	}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to insert text: %v", err)
	}

	return nil
}

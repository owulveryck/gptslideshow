package main

import (
	"context"
	"fmt"

	slides "google.golang.org/api/slides/v1"
)

// AddSpeakerNote adds a speaker note to a specific slide in a Google Slides presentation.
func AddSpeakerNote(srv *slides.Service, presentationID, slideID, noteContent string) error {
	// Create the request to update speaker notes.
	requests := []*slides.Request{
		{
			UpdateSlideProperties: &slides.UpdateSlidePropertiesRequest{
				ObjectId: slideID,
				Fields:   "notesPage",
				SlideProperties: &slides.SlideProperties{
					NotesPage: &slides.Page{
						NotesProperties: &slides.NotesProperties{
							SpeakerNotesObjectId: slideID,
						},
						PageElements: []*slides.PageElement{
							{
								Shape: &slides.Shape{
									Text: &slides.TextContent{
										TextElements: []*slides.TextElement{
											{
												TextRun: &slides.TextRun{
													Content: noteContent,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Execute the batch update request.
	_, err := srv.Presentations.BatchUpdate(presentationID, &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}).Context(context.Background()).Do()
	if err != nil {
		return fmt.Errorf("unable to add speaker note: %v", err)
	}

	return nil
}

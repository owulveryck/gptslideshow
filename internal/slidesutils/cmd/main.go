package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/owulveryck/gptslideshow/internal/gcputils"
	"github.com/owulveryck/gptslideshow/internal/slidesutils"
	"github.com/owulveryck/gptslideshow/internal/slidesutils/mytemplate"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/slides/v1"
)

func main() {
	var presentationID string

	flag.StringVar(&presentationID, "id", "", "ID of the slide to update, empty means create a new one")
	flag.Parse()
	ctx := context.Background()

	// Load client secret file
	b, err := os.ReadFile("../../../../credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Initialize Google OAuth2 config
	config, err := google.ConfigFromJSON(b, drive.DriveScope, slides.PresentationsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Get authenticated client
	client := gcputils.GetClient(config)

	// Initialize Google Slides service
	slidesSrv, err := slides.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}

	// Initialize the Builder
	builder, err := mytemplate.NewBuilder(ctx, slidesSrv, presentationID)
	if err != nil {
		log.Fatalf("Unable to create Builder: %v", err)
	}
	err = builder.CreateNewSlide(ctx, mytemplate.TitleSubtitleBody)
	// Find placeholders for title, subtitle, and body in the newly created slide.
	var titlePlaceholderID, subtitlePlaceholderID, bodyPlaceholderID string
	for _, element := range builder.CurrentSlide.PageElements {
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
	_ = titlePlaceholderID
	_ = subtitlePlaceholderID

	bold := slidesutils.EncodeStyle(true, false)
	italic := slidesutils.EncodeStyle(false, true)
	normal := slidesutils.EncodeStyle(false, false)
	boldItalic := slidesutils.EncodeStyle(true, true)
	content := []slidesutils.Chunk{
		{Content: "this is a ", Style: normal, IndentationLevel: 0},
		{Content: "bold", Style: bold, IndentationLevel: 0},
		{Content: " word and this is an ", Style: normal, IndentationLevel: 0},
		{Content: "italic", Style: italic, IndentationLevel: 0},
		{Content: " like ", Style: normal, IndentationLevel: 0},
		{Content: "this", Style: boldItalic, IndentationLevel: 0, Paragraph: 0},
		{Content: "This is a list:", Style: normal, IndentationLevel: 0, Paragraph: 1},
		{Content: "the level of indentation should be ", Style: normal, IndentationLevel: 1, Paragraph: 2},
		{Content: "1", Style: bold, IndentationLevel: 1, Paragraph: 2},
		{Content: "the level of indentation should be ", Style: normal, IndentationLevel: 1, Paragraph: 3},
		{Content: "1", Style: bold, IndentationLevel: 1, Paragraph: 3},
		{Content: "this Content should have a level indentation of 2", Style: normal, IndentationLevel: 2, Paragraph: 4},
		{Content: "the level of indentation should be ", Style: normal, IndentationLevel: 1, Paragraph: 5},
		{Content: "1", Style: bold, IndentationLevel: 1, Paragraph: 5},
		{Content: "and this is back to a level of indentation of zero", Style: normal, IndentationLevel: 0, Paragraph: 6},
		{Content: "and this is back to a level of indentation of zero", Style: bold, IndentationLevel: 0, Paragraph: 7},
	}
	requests := slidesutils.InsertMarkdownContent(content, bodyPlaceholderID)
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(requests)

	// Create the BatchUpdate request.
	batchUpdateRequest := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}

	// Call the Slides API to execute the batch update.
	_, err = builder.Srv.Presentations.BatchUpdate(presentationID, batchUpdateRequest).Do()
	if err != nil {
		log.Fatalf("Unable to update slide: %v", err)
	}
}

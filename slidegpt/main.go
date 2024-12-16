package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/owulveryck/gptslideshow/config"
	"github.com/owulveryck/gptslideshow/internal/ai"
	"github.com/owulveryck/gptslideshow/internal/slidesutils"
	"github.com/owulveryck/gptslideshow/internal/slidesutils/mytemplate"
	"github.com/owulveryck/gptslideshow/internal/structure"
	slides "google.golang.org/api/slides/v1"
)

func main() {
	// Parse command-line flags
	presentationId := flag.String("id", "", "ID of the slide to update, empty means create a new one")

	flag.Parse()
	ctx := context.Background()
	openaiClient := ai.NewAI()

	// Initialize Google services
	client := initGoogleClient()
	slidesSrv := initSlidesService(client)

	builder, err := mytemplate.NewBuilder(ctx, slidesSrv, *presentationId)
	if err != nil {
		log.Fatal(err)
	}
	builder.CreateSlideTitleSubtitleBody(ctx, structure.Slide{
		Title:    "SlideGPT",
		Subtitle: "Waiting for instruction, write GoGoGo to send",
		Body:     "",
		Chapter:  false,
	})
	found := false
	var content string
	var bodyID string
	for !found {
		//	currentSlide, _ := getSlideByID(builder.Presentation.PresentationId, builder.CurrentSlide.ObjectId, builder.Srv)
		for _, currentSlide := range builder.Presentation.Slides {
			if found {
				break
			}
			// Find the BODY placeholder
			for _, element := range currentSlide.PageElements {
				if found {
					break
				}
				var currentText string
				if element.Shape != nil && element.Shape.Placeholder != nil {
					bodyID = element.ObjectId
					// Check if Text and TextElements exist
					if element.Shape != nil && element.Shape.Text != nil && len(element.Shape.Text.TextElements) > 0 {
						for i := range element.Shape.Text.TextElements {
							textRun := element.Shape.Text.TextElements[i].TextRun
							if textRun != nil {
								currentText += textRun.Content
								fmt.Printf("Text: %s\n", currentText)
								if strings.Contains(currentText, "CHATGPT:") {
									content = strings.Replace(currentText, "CHATGPT:", "", -1)
									found = true
									break
								}
							}
						}
					} else {
						fmt.Println("No text found in BODY placeholder.")
					}
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
	textRequests := []*slides.Request{
		{
			DeleteText: &slides.DeleteTextRequest{
				ObjectId: bodyID,
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       bodyID,
				InsertionIndex: 0,
				Text:           "Sending to ChatGPT...\n\n",
			},
		},
	}
	//
	if _, err := builder.Srv.Presentations.BatchUpdate(builder.Presentation.PresentationId, &slides.BatchUpdatePresentationRequest{
		Requests: textRequests,
	}).Context(ctx).Do(); err != nil {
		log.Fatal("failed to insert text: %w", err)
	}

	// Query OpenAI API for validation or enhancement (optional)
	chat, err := openaiClient.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(content),
		}),
		Model: openai.F(config.ConfigInstance.OpenAIModel),
		// Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})
	if err != nil {
		log.Fatal(err)
	}
	res := []*slides.Request{
		{
			DeleteText: &slides.DeleteTextRequest{
				ObjectId: bodyID,
			},
		},
	}

	resChat := slidesutils.InsertMarkdownContent(chat.Choices[0].Message.Content, bodyID)
	// Sort the requests slice
	sort.SliceStable(resChat, func(i, j int) bool {
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
		return priority(resChat[i]) < priority(resChat[j])
	})
	log.Println(chat.Choices[0].Message.Content)

	res = append(res, resChat...)

	for _, res := range res {
		if _, err := builder.Srv.Presentations.BatchUpdate(builder.Presentation.PresentationId, &slides.BatchUpdatePresentationRequest{
			Requests: []*slides.Request{res},
		}).Context(ctx).Do(); err != nil {
			log.Fatal("failed to insert text: %w", err)
		}

		time.Sleep(10 * time.Millisecond)
	}
	_ = openaiClient
}

func getSlideByID(presentationID, slideID string, service *slides.Service) (*slides.Page, error) {
	// Retrieve the presentation
	presentation, err := service.Presentations.Get(presentationID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve presentation: %w", err)
	}

	// Iterate through slides to find the one with the matching ID
	for _, slide := range presentation.Slides {
		if slide.ObjectId == slideID {
			return slide, nil
		}
	}

	// Slide not found
	return nil, fmt.Errorf("slide with ID %s not found", slideID)
}

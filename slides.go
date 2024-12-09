package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	drive "google.golang.org/api/drive/v3"

	"github.com/owulveryck/gptslideshow/internal/ai"
	"github.com/owulveryck/gptslideshow/internal/driveutils"
	"github.com/owulveryck/gptslideshow/internal/slidesutils"
	"github.com/owulveryck/gptslideshow/internal/structure"
)

func generateSlides(ctx context.Context, openaiClient *ai.AI, prompt string, content []byte) *structure.Presentation {
	saveContent("prompt-*.txt", []byte(prompt))
	presentationData, err := openaiClient.GenerateContentFromText(ctx, prompt, content)
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.MarshalIndent(presentationData, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	saveContent("generated-data-*.json", b)

	return presentationData
}

func createPresentationSlides(ctx context.Context, builder slidesutils.BuilderInterface, driveSrv *drive.Service, openaiClient *ai.AI, withImages bool, presentationData *structure.Presentation) error {
	err := builder.CreateCover(ctx, presentationData.Title, presentationData.Subtitle)
	if err != nil {
		return err
	}

	for i, slide := range presentationData.Slides {
		log.Printf("Slide %v: %v", i, slide.Title)
		if slide.Chapter {
			err = builder.CreateChapter(ctx, slide)
			if err != nil {
				return err
			}
			if withImages {
				// Generate the illustration
				img, err := openaiClient.GenerateImageFromText(ctx, slide.Body)
				if err != nil {
					return err
				}
				imageUrl, err := driveutils.UploadImage(ctx, driveSrv, img, slide.Title+".png")
				if err != nil {
					return err
				}
				var imageWidth, imageHeight, slideWidth, slideHeight, translateX, translateY float64
				// Dimensions in EMUs (e.g., 3x3 inches)
				imageWidth = 2743200  // 3 inches
				imageHeight = 2743200 // 3 inches

				// Slide dimensions in EMUs
				slideWidth = 9144000  // 10 inches
				slideHeight = 6858000 // 7.5 inches

				// Calculate position to center the image
				translateX = (slideWidth - imageWidth) / 2
				translateY = (slideHeight - imageHeight) / 2
				translateX = 1213950
				translateY = 1659800

				// Insert the image centered on the current slide
				err = builder.InsertImage(ctx, imageUrl, imageWidth, imageHeight, translateX, translateY)
				if err != nil {
					return err
				}
			}
		} else {
			err = builder.CreateSlideTitleSubtitleBody(ctx, slide)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("New presentation created and modified successfully.")
	return nil
}

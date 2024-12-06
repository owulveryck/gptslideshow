package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	drive "google.golang.org/api/drive/v3"
	slides "google.golang.org/api/slides/v1"

	"github.com/owulveryck/gptslideshow/internal/ai"
	"github.com/owulveryck/gptslideshow/internal/driveutils"
	"github.com/owulveryck/gptslideshow/internal/slidesutils"
	"github.com/owulveryck/gptslideshow/internal/structure"
)

func generateSlides(ctx context.Context, openaiClient *ai.AI, prompt string, content []byte) *structure.Presentation {
	presentationData, err := openaiClient.GenerateContentFromText(ctx, prompt, content)
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(presentationData)

	return presentationData
}

func createPresentationSlides(ctx context.Context, slidesSrv *slides.Service, driveSrv *drive.Service, openaiClient *ai.AI, withImages bool, presentationId string, presentationData *structure.Presentation) error {
	builder, err := slidesutils.NewBuilder(ctx, slidesSrv, presentationId)
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

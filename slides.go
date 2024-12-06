package main

import (
	"context"
	"fmt"
	"log"

	"github.com/owulveryck/gptslideshow/internal/ai"
	"github.com/owulveryck/gptslideshow/internal/slidesutils"
	"github.com/owulveryck/gptslideshow/internal/structure"
	slides "google.golang.org/api/slides/v1"
)

func generateSlides(ctx context.Context, openaiClient *ai.AI, prompt string, content []byte) *structure.Presentation {
	presentationData, err := ai.NewAI().GenerateContentFromText(ctx, prompt, content)
	if err != nil {
		log.Fatal(err)
	}
	return presentationData
}

func createPresentationSlides(ctx context.Context, slidesSrv *slides.Service, presentationId string, presentationData *structure.Presentation) {
	builder, err := slidesutils.NewBuilder(ctx, slidesSrv, presentationId)
	if err != nil {
		log.Fatalf("Unable to create Builder: %v", err)
	}

	for i, slide := range presentationData.Slides {
		log.Printf("Slide %v: %v", i, slide.Title)
		if slide.Chapter {
			err = builder.CreateChapter(ctx, slide)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = builder.CreateSlideTitleSubtitleBody(ctx, slide)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("New presentation created and modified successfully.")
}

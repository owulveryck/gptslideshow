package main

import (
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/api/slides/v1"
)

func ListSlidesAndElements(srv *slides.Service, presentationId string) error {
	// Retrieve the presentation
	presentation, err := srv.Presentations.Get(presentationId).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve presentation: %v", err)
	}

	fmt.Printf("Presentation Title: %s\n", presentation.Title)
	fmt.Println("Listing slides and elements:")
	layoutMap := make(map[string]string)
	for _, layout := range presentation.Layouts {
		layoutMap[layout.ObjectId] = layout.LayoutProperties.Name
	}
	for k, v := range layoutMap {
		fmt.Printf("%v : %v\n", k, v)
	}

	// Iterate through the slides
	for i, slide := range presentation.Slides {
		fmt.Printf("\nSlide %d:\n", i+1)

		// Retrieve layout name
		layoutId := slide.SlideProperties.LayoutObjectId
		fmt.Printf("  Layout Name: %s\n", layoutId)

		// List elements on the slide
		fmt.Println("  Elements:")
		for _, element := range slide.PageElements {
			content, err := json.MarshalIndent(element, "", "  ")
			if err != nil {
				log.Fatalf("Error marshaling JSON: %v", err)
			}
			fmt.Println(string(content))
			/*
				if element.Shape != nil {
					elementType := element.Shape.ShapeType
					elementId := element.ObjectId
					var elementText []byte
					if element.Shape.Text != nil {
						// Assuming element.Shape.Text implements the `json.Marshaler` interface
						elementText, err = json.MarshalIndent(element.Shape.Text, "", "  ")
						if err != nil {
							log.Fatalf("Error marshaling JSON: %v", err)
						}
					}
					fmt.Printf("    - ID: %s, Type: %s, Text: %s\n", elementId, elementType, elementText)
				}
			*/
		}
	}

	return nil
}

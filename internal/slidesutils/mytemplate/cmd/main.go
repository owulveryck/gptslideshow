package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/owulveryck/gptslideshow/internal/gcputils"
	"github.com/owulveryck/gptslideshow/internal/slidesutils/mytemplate"
	"github.com/owulveryck/gptslideshow/internal/structure"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
	slides "google.golang.org/api/slides/v1"
)

func main() {
	var presentationId string

	flag.StringVar(&presentationId, "id", "", "ID of the slide to update, empty means create a new one")
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
	builder, err := mytemplate.NewBuilder(ctx, slidesSrv, presentationId)
	if err != nil {
		log.Fatalf("Unable to create Builder: %v", err)
	}

	log.Println("Generating slides...")
	slide := structure.Slide{
		Title:    "Title of the slide",
		Subtitle: "Subtitle of the slide",
		Body: `this is a **bold** word and this is a list: Hello, 世界\xF0\x28\x8C\x28InvalidUTF8!
- the level of indentation should be 1
  - this content should have a level indentation of 2
and this is back to a level of indentation of zero` + "Inspired by Simon Wardley's theory of evolution, this slide categorizes data progression in organizations into four phases:\n\n1. **Genesis** – Where data is initially unstructured, akin to the experimental world of startups.\n2. **Craft** – Early stages of structuring data and building robust applications.\n3. **Product** – Data is now an asset, structured for efficient usability across departments.\n4. **Commodity** – Data becomes a ubiquitous part of the organizational ecosystem, similar to basic utilities in society.\n\nUsing a two-dimensional model, we map the omnipresence and certainty of data over time, illustrating a diffusion curve similar to that seen in technology adoption cycles.",
		Chapter: false,
	}

	err = builder.CreateCover(ctx, "AA", "BB")
	if err != nil {
		log.Fatal(err)
	}
	// Create a slide with title, subtitle, and body
	//	err = builder.CreateSlideTitleSubtitleBody(ctx, slide)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	// Create a chapter slide
	err = builder.CreateChapter(ctx, slide)
	if err != nil {
		log.Fatal(err)
	}
	// Define the image properties
	imageUrl := "https://upload.wikimedia.org/wikipedia/commons/8/8d/Sinclair_QL_256x256_mode_example_image.png"

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
		log.Fatalf("Error inserting image: %v", err)
	}

	log.Println("Image inserted successfully.")

	fmt.Println("New presentation created and modified successfully.")

	// Example Markdown content
	markdown := `

This is a paragraph with some **bold** text.

- First item
- Second item
- Third item
`

	// Parse the Markdown into an AST
	md := goldmark.New()
	source := []byte(markdown)
	reader := text.NewReader(source)
	parser := md.Parser()
	doc := parser.Parse(reader)
	err = builder.CreateNewSlide(ctx, mytemplate.TitleSubtitleBody)

	// Add the AST content to the current slide
	err = builder.AddNodeContent(doc, source)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

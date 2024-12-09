package main

import (
	"flag"
	"fmt"

	"github.com/owulveryck/gptslideshow/config"
)

func parseFlags() (presentationId *string, fromTemplate *string, prompt *string, textfile *string, audiofile *string, helpFlag *bool) {
	presentationId = flag.String("id", "", "ID of the slide to update, empty means create a new one")
	fromTemplate = flag.String("t", "", "ID of a template file")
	helpFlag = flag.Bool("h", false, "help")
	prompt = flag.String("prompt", `Convert the following text into an array of structured slides.
Each slide should have a title, a subtitle, and a body that should add comprehensive and detailed explanation. Do not use markdown, and seperate each paragraph with two newlines;

You can also generate chapters between a set of content slides.
If the slide is a chapter, the body should contain a complete description of the content of the chapter usable to generate a picture to illustrate.

If it is a chapter, the field 'chapter' must be set to true.
		The first slide should be an executive summary. Generate the most complete possible output. Here is the content:


`, "the prompt")

	textfile = flag.String("content", "", "The content file")
	audiofile = flag.String("audio", "", "The audio file in mp3")

	flag.Parse()
	return
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  [flags]")

	fmt.Println("\nFlags:")
	flag.PrintDefaults()

	fmt.Println("\nEnvironment Variables:")
	config.Help()
}

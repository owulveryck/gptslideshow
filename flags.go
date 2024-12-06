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
	prompt = flag.String("prompt", `Convert the following text into an array of structured slides...`, "The prompt")
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

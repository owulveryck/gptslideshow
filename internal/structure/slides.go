package structure

import "github.com/invopop/jsonschema"

// Presentation represents the entire presentation structure
type Presentation struct {
	OriginalContent []byte  `json:"-"`
	Title           string  `json:"presentation_title" jsonschema_description:"The title of the presentation"`
	Subtitle        string  `json:"presentation_subtitle" jsonschema_description:"The subtitle of the presentation"`
	Slides          []Slide `json:"slides" jsonschema_description:"The content of the presentation"`
}

// Slide represents a single slide in the presentation
type Slide struct {
	Title        string       `json:"title" jsonschema_description:"The title of the slide"`
	Subtitle     string       `json:"subtitle" jsonschema_description:"The subtitle of the slide"`
	Body         string       `json:"body" jsonschema_description:"The main content of the slide or the description of the chapter"`
	Chapter      bool         `json:"chapter" jsonschema_description:"A boolean to indicate if this slides introduces a new chapter"`
	Illustration Illustration `json:"illustration" jsonschema_description:"A nested structure to generate an illustration if relevant and only if relevant to the current slide. If not leave it blank"`
}

type Illustration struct {
	Prompt string `json:"prompt" jsonschema_description:"A prompt to generate the illustration"`
	Alt    string `json:"alt" jsonschema_description:"The alternate text of the illustration"`
}

// GenerateSchema generates the JSON schema for a given type
func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	return reflector.Reflect(v)
}

// Generate the JSON schema for Slide
var (
	PresentationResponseSchema = GenerateSchema[Presentation]()
	SlideResponseSchema        = GenerateSchema[Slide]()
)

package slidesutils

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"google.golang.org/api/slides/v1"
)

type chunk struct {
	content          string
	isBold           bool
	indentationLevel int // 0 means no indentation, 1 first bullet, 2 second bulled
}

func Format(content string, objectID string) []*slides.Request {
	chunks := parseContent(content)

	// Build the text content with styles
	var requests []*slides.Request
	currentIndex := int64(0) // Tracks the cumulative character index in the text box

	for _, c := range chunks {

		// Append the text to the text box
		requests = append(requests, &slides.Request{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       objectID,
				InsertionIndex: currentIndex,
				Text:           c.content, // Add a newline after each chunk
			},
		})

		// Calculate the start and end indices for the current chunk
		startIndex := currentIndex
		endIndex := startIndex + int64(utf8.RuneCountInString(c.content))
		// Replace only if the string ends with "\n\n"
		//		if strings.HasSuffix(c.content, "\n\n") {
		//			endIndex--
		//		}

		// Apply bold styling if needed
		if c.isBold {
			requests = append(requests, &slides.Request{
				UpdateTextStyle: &slides.UpdateTextStyleRequest{
					ObjectId: objectID,
					TextRange: &slides.Range{
						Type:       "FIXED_RANGE",
						StartIndex: &startIndex,
						EndIndex:   &endIndex,
					},
					Style: &slides.TextStyle{
						Bold: true,
					},
					Fields: "bold",
				},
			})
		} else {
			requests = append(requests, &slides.Request{
				UpdateTextStyle: &slides.UpdateTextStyleRequest{
					ObjectId: objectID,
					TextRange: &slides.Range{
						Type:       "FIXED_RANGE",
						StartIndex: &startIndex,
						EndIndex:   &endIndex,
					},
					Style: &slides.TextStyle{
						Bold: false,
					},
					Fields: "bold",
				},
			})
		}

		// Apply bullet points based on indentation level
		if c.indentationLevel > 0 {
			// Create bullet points
			requests = append(requests, &slides.Request{
				CreateParagraphBullets: &slides.CreateParagraphBulletsRequest{
					ObjectId: objectID,
					TextRange: &slides.Range{
						Type:       "FIXED_RANGE",
						StartIndex: &startIndex,
						EndIndex:   &endIndex, // Include the newline for bullet points
					},
					BulletPreset: "BULLET_DISC_CIRCLE_SQUARE",
				},
			})

			// Adjust indentation level
			var indentStart, indentFirstStart float64
			switch c.indentationLevel {
			case 1:
				indentFirstStart = 18.0 // First level indentation, in points
				indentStart = 36.0      // First level indentation, in points
			case 2:
				indentFirstStart = 54.0 // First level indentation, in points
				indentStart = 72.0      // Second level indentation, in points
			default:
				indentStart = float64(c.indentationLevel) * 36.0 // Adjust as needed for deeper levels
			}

			requests = append(requests, &slides.Request{
				UpdateParagraphStyle: &slides.UpdateParagraphStyleRequest{
					ObjectId: objectID,
					TextRange: &slides.Range{
						Type:       "FIXED_RANGE",
						StartIndex: &startIndex,
						EndIndex:   &endIndex,
					},
					Style: &slides.ParagraphStyle{
						IndentFirstLine: &slides.Dimension{
							Magnitude: indentFirstStart,
							Unit:      "PT",
						},
						IndentStart: &slides.Dimension{
							Magnitude: indentStart,
							Unit:      "PT",
						},
					},
					Fields: "indentStart, indentFirstLine", // Ensure only the indentStart field is updated
				},
			})
		}

		// Update the current index to account for the inserted text and newline
		currentIndex = endIndex
	}

	return requests
}

func parseContent(input string) []chunk {
	var chunks []chunk

	lines := strings.Split(input, "\n")              // Split the input into lines
	boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*`) // Regex to detect bold text

	for _, line := range lines {
		trimmedLine := filterPrintable(line)

		// Determine the indentation level
		var indentationLevel int
		if strings.HasPrefix(trimmedLine, "- ") {
			indentationLevel = 1
			trimmedLine = strings.TrimPrefix(trimmedLine, "- ")
		} else if strings.HasPrefix(trimmedLine, "  - ") {
			indentationLevel = 2
			trimmedLine = strings.TrimPrefix(trimmedLine, "  - ")
		}

		// Split the line into chunks by bold markers
		parts := boldRegex.Split(trimmedLine, -1)
		matches := boldRegex.FindAllStringSubmatch(trimmedLine, -1)

		// Iterate over the parts and matches to build the chunks
		for i, part := range parts {
			if part != "" {
				chunks = append(chunks, chunk{
					content:          part,
					isBold:           false,
					indentationLevel: indentationLevel,
				})
			}

			// Add the bold content if there's a match
			if i < len(matches) {
				chunks = append(chunks, chunk{
					content:          matches[i][1], // The content inside the bold markers
					isBold:           true,
					indentationLevel: indentationLevel,
				})
			}
		}
		lastChunk := chunks[len(chunks)-1]
		lastChunk.content = lastChunk.content + "\n"
		chunks[len(chunks)-1] = lastChunk
	}
	// Remove the last carriage rerturn
	lastChunk := chunks[len(chunks)-1]
	lastChunk.content = lastChunk.content[:len(lastChunk.content)-1]
	chunks[len(chunks)-1] = lastChunk

	return chunks
}

func filterPrintable(s string) string {
	result := []rune{}
	for _, r := range s {
		if unicode.IsSpace(r) {
			result = append(result, ' ')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

package slidesutils

import (
	"log"
	"unicode/utf8"

	"google.golang.org/api/slides/v1"
)

func Format(content string, objectID string) []*slides.Request {
	currentParagraph := 0
	chunks := parseContent(content)

	// Build the text content with styles
	var requests []*slides.Request
	currentIndex := int64(0) // Tracks the cumulative character index in the text box

	for _, c := range chunks {
		cr := ""
		if c.paragraph != currentParagraph {
			cr = "\n"
			currentIndex = int64(c.paragraph)
		}
		log.Printf("content: '%v'", c.content)

		// Append the text to the text box
		requests = append(requests, &slides.Request{
			InsertText: &slides.InsertTextRequest{
				ObjectId:       objectID,
				InsertionIndex: currentIndex,
				Text:           c.content + cr, // Add a newline after each chunk
			},
		})

		// Calculate the start and end indices for the current chunk
		startIndex := currentIndex
		endIndex := startIndex + int64(utf8.RuneCountInString(c.content+cr))
		bold, italic, normal := decodeStyle(c.style)

		switch {
		case bold:
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
		case normal:
			requests = append(requests, &slides.Request{
				UpdateTextStyle: &slides.UpdateTextStyleRequest{
					ObjectId: objectID,
					TextRange: &slides.Range{
						Type:       "FIXED_RANGE",
						StartIndex: &startIndex,
						EndIndex:   &endIndex,
					},
					Style: &slides.TextStyle{
						Italic: false,
						Bold:   false,
					},
					Fields: "*",
				},
			})
		case italic:
			requests = append(requests, &slides.Request{
				UpdateTextStyle: &slides.UpdateTextStyleRequest{
					ObjectId: objectID,
					TextRange: &slides.Range{
						Type:       "FIXED_RANGE",
						StartIndex: &startIndex,
						EndIndex:   &endIndex,
					},
					Style: &slides.TextStyle{
						Italic: true,
					},
					Fields: "*",
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

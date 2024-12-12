package slidesutils

import (
	"unicode/utf8"

	"google.golang.org/api/slides/v1"
)

// InsertText inserts text into a specified object with a given insertion index.
func InsertText(objectID string, content string, insertionIndex int64) *slides.Request {
	return &slides.Request{
		InsertText: &slides.InsertTextRequest{
			ObjectId:       objectID,
			InsertionIndex: insertionIndex,
			Text:           content,
		},
	}
}

// UpdateTextStyle updates the text style (e.g., bold, italic) for a specific range.
func UpdateTextStyle(objectID string, startIndex, endIndex int64, bold, italic, normal bool) *slides.Request {
	return &slides.Request{
		UpdateTextStyle: &slides.UpdateTextStyleRequest{
			ObjectId: objectID,
			TextRange: &slides.Range{
				Type:       "FIXED_RANGE",
				StartIndex: &startIndex,
				EndIndex:   &endIndex,
			},
			Style: &slides.TextStyle{
				Bold:   bold,
				Italic: italic,
			},
			Fields: "bold,italic",
		},
	}
}

// CreateBullets creates bullet points for a specific range.
func CreateBullets(objectID string, startIndex, endIndex int64) *slides.Request {
	return &slides.Request{
		CreateParagraphBullets: &slides.CreateParagraphBulletsRequest{
			ObjectId: objectID,
			TextRange: &slides.Range{
				Type:       "FIXED_RANGE",
				StartIndex: &startIndex,
				EndIndex:   &endIndex,
			},
			BulletPreset: "BULLET_DISC_CIRCLE_SQUARE",
		},
	}
}

// UpdateParagraphStyle updates the paragraph style (e.g., indentation) for a specific range.
func UpdateParagraphStyle(objectID string, startIndex, endIndex int64, indentFirstLine, indentStart float64) *slides.Request {
	return &slides.Request{
		UpdateParagraphStyle: &slides.UpdateParagraphStyleRequest{
			ObjectId: objectID,
			TextRange: &slides.Range{
				Type:       "FIXED_RANGE",
				StartIndex: &startIndex,
				EndIndex:   &endIndex,
			},
			Style: &slides.ParagraphStyle{
				IndentFirstLine: &slides.Dimension{
					Magnitude: indentFirstLine,
					Unit:      "PT",
				},
				IndentStart: &slides.Dimension{
					Magnitude: indentStart,
					Unit:      "PT",
				},
			},
			Fields: "indentFirstLine,indentStart",
		},
	}
}

// InsertMarkdownContent processes the content and generates a list of requests for formatting.
func InsertMarkdownContent(input string, objectID string) []*slides.Request {
	chunks := parseContent(input)
	// toRemove is the number of tabs to remove from the endIndex
	toRemove := int64(0)
	var requests []*slides.Request
	currentIndex := int64(0) // Tracks the cumulative character index in the text box.

	inList := false
	var inListStartIndex, inListEndIndex int64
	// First, insert all the text
	for i, c := range chunks {
		if i < len(chunks)-1 && chunks[i+1].Paragraph != c.Paragraph {
			c.Content += "\n"
		}

		// Insert text into the text box.
		if i > 0 {
			if chunks[i-1].Paragraph != c.Paragraph {
				if c.IndentationLevel == 1 {
					c.Content = "" + c.Content
				}
				if c.IndentationLevel == 2 {
					c.Content = "\t" + c.Content
					toRemove++
				}
			}
		}
		startIndex := currentIndex
		endIndex := startIndex + int64(utf8.RuneCountInString(c.Content))

		requests = append(requests, InsertText(objectID, c.Content, currentIndex))

		// Decode and apply text styles.
		bold, italic, normal := DecodeStyle(c.Style)
		if bold || italic || normal {
			requests = append(requests, UpdateTextStyle(objectID, startIndex, endIndex, bold, italic, normal))
		}
		// We are entering a list
		if c.IndentationLevel > 0 && !inList {
			inListStartIndex = startIndex
			inList = true
		}
		// Update the current index to account for the inserted text and newline.
		currentIndex = endIndex
		// End of the list Apply Request
		if (c.IndentationLevel == 0 || i == len(chunks)-1) && inList {
			requests = append(requests, &slides.Request{
				CreateParagraphBullets: &slides.CreateParagraphBulletsRequest{
					BulletPreset: "BULLET_DISC_CIRCLE_SQUARE",
					ObjectId:     objectID,
					TextRange: &slides.Range{
						StartIndex: &inListStartIndex,
						EndIndex:   &inListEndIndex,
						Type:       "FIXED_RANGE",
					},
				},
			})
			inList = false
			currentIndex -= toRemove
			toRemove = 0
		}

		if inList {
			inListEndIndex = endIndex - toRemove
		}
	}

	return requests
}

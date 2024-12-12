package mytemplate

import (
	"bytes"
	"fmt"

	"github.com/owulveryck/gptslideshow/internal/slidesutils"

	"github.com/yuin/goldmark/ast"
	slides "google.golang.org/api/slides/v1"
)

// AddNodeContent processes an AST node and adds its content to the current slide.
func (b *Builder) AddNodeContent(node ast.Node, source []byte) error {
	var titleID, bodyID string
	for _, element := range b.CurrentSlide.PageElements {
		if element.Shape != nil && element.Shape.Placeholder != nil {
			switch element.Shape.Placeholder.Type {
			case "TITLE":
				titleID = element.ObjectId
			case "BODY":
				bodyID = element.ObjectId
			}
		}
	}

	// Ensure that CurrentSlide is set
	if b.CurrentSlide == nil {
		return fmt.Errorf("no current slide is set")
	}

	// Recursive function to traverse the AST and add content to the slide
	var traverse func(node ast.Node) error
	traverse = func(node ast.Node) error {
		switch node.Kind() {
		case ast.KindHeading:
			heading := node.(*ast.Heading)
			text := extractText(heading, source)
			if heading.Level == 1 {
				// Add H1 headings to the title placeholder
				err := b.AddTextToPlaceholder(text, titleID)
				if err != nil {
					return err
				}
			} else {
				// Add other headings to the body placeholder
				err := b.AddTextToPlaceholder(fmt.Sprintf("%s\n", text), bodyID)
				if err != nil {
					return err
				}
			}
		case ast.KindParagraph:
			paragraph := node.(*ast.Paragraph)
			text := extractText(paragraph, source)
			err := b.AddTextToPlaceholder(fmt.Sprintf("%s\n", text), bodyID)
			if err != nil {
				return err
			}
		case ast.KindList:
			list := node.(*ast.List)
			err := b.processList(list, bodyID, source)
			if err != nil {
				return err
			}
		}

		// Recursively process child nodes
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			if err := traverse(child); err != nil {
				return err
			}
		}
		return nil
	}

	return traverse(node)
}

// AddTextToPlaceholder writes text into a placeholder of the current slide using slidesutils.
func (b *Builder) AddTextToPlaceholder(text, objectID string) error {
	if b.CurrentSlide == nil {
		return fmt.Errorf("no current slide to add text to")
	}
	// Use slidesutils.InsertText to create the request
	request := slidesutils.InsertText(objectID, text, 0)

	// Send the batch update request
	batchUpdateRequest := &slides.BatchUpdatePresentationRequest{
		Requests: []*slides.Request{request},
	}
	_, err := b.Srv.Presentations.BatchUpdate(b.Presentation.PresentationId, batchUpdateRequest).Do()
	if err != nil {
		return fmt.Errorf("failed to add text to placeholder: %v", err)
	}
	return nil
}

// processList processes a list node and adds it as a bulleted list to the current slide using slidesutils.
func (b *Builder) processList(list *ast.List, objectID string, source []byte) error {
	if b.CurrentSlide == nil {
		return fmt.Errorf("no current slide to add list to")
	}

	// Combine list items into a single string, separated by newlines
	var buf bytes.Buffer
	for item := list.FirstChild(); item != nil; item = item.NextSibling() {
		buf.WriteString(extractText(item, source))
		buf.WriteString("\n")
	}

	// Add the combined list text to the body placeholder
	text := buf.String()
	err := b.AddTextToPlaceholder(text, objectID)
	if err != nil {
		return fmt.Errorf("failed to add list text to placeholder: %v", err)
	}

	// Create bullets for the list using slidesutils.CreateBullets
	request := slidesutils.CreateBullets(objectID, 0, int64(len(text)))
	batchUpdateRequest := &slides.BatchUpdatePresentationRequest{
		Requests: []*slides.Request{request},
	}
	_, err = b.Srv.Presentations.BatchUpdate(b.Presentation.PresentationId, batchUpdateRequest).Do()
	if err != nil {
		return fmt.Errorf("failed to create bullets for list: %v", err)
	}

	return nil
}

// extractText recursively extracts plain text from an AST node and its children, using the source content.
func extractText(node ast.Node, source []byte) string {
	var buf bytes.Buffer

	// Walk through all children of the node
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch child := child.(type) {
		case *ast.Text:
			// Use the source content to extract the segment's value
			buf.Write(child.Segment.Value(source))

			// Check if the text node has a hard line break
			if child.HardLineBreak() {
				buf.WriteString("\n")
			}
		default:
			// If the child is not a text node, recursively extract its text
			buf.WriteString(extractText(child, source))
		}
	}

	return buf.String()
}

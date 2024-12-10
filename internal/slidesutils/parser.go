package slidesutils

import (
	"log"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Define the chunk struct
type chunk struct {
	content          string
	style            style
	indentationLevel int // 0 means no indentation, 1 first bullet, 2 second bullet
	paragraph        int
}

// Define the custom type
type style []byte

// Define bit masks for each style
const (
	normalMask byte = 1 << 0 // 0000 0001
	boldMask   byte = 1 << 1 // 0000 0010
	italicMask byte = 1 << 2 // 0000 0100
)

// Encode a style into a custom style type
func encodeStyle(bold, italic bool) style {
	var s byte
	if bold {
		s |= boldMask
	}
	if italic {
		s |= italicMask
	}
	if !italic && !bold {
		s |= normalMask
	}
	return style{s}
}

// Decode a style back into its components
func decodeStyle(s style) (bold, italic, normal bool) {
	if len(s) == 0 {
		return false, false, false
	}
	bold = s[0]&boldMask != 0
	italic = s[0]&italicMask != 0
	normal = s[0]&normalMask != 0
	return
}

var paragraph int

// Recursive function to traverse AST and collect chunks
func processNode(n ast.Node, reader text.Reader, level int, currentStyle style, chunks *[]chunk) {
	switch n.Kind() {
	case ast.KindEmphasis:
	case ast.KindText:
	case ast.KindList:
	case ast.KindListItem:
		paragraph++
	case ast.KindParagraph:
		paragraph++
	case ast.KindTextBlock:
	case ast.KindDocument:
	default:
		log.Printf("Warning, node kind %v will may not be rendered correctly", n.Kind().String())
	}
	if textNode, ok := n.(*ast.Text); ok {
		// Extract text content
		content := textNode.Segment.Value(reader.Source())
		*chunks = append(*chunks, chunk{
			content:          string(content),
			style:            currentStyle,
			indentationLevel: level,
			paragraph:        paragraph,
		})
	}

	if emphasisNode, ok := n.(*ast.Emphasis); ok {
		// Update style for emphasis
		newStyle := currentStyle
		bold, italic, _ := decodeStyle(currentStyle)
		if emphasisNode.Level == 1 {
			italic = true
		} else if emphasisNode.Level == 2 {
			bold = true
		}
		newStyle = encodeStyle(bold, italic)

		// Process children with updated style
		for child := emphasisNode.FirstChild(); child != nil; child = child.NextSibling() {
			processNode(child, reader, level, newStyle, chunks)
		}
		return
	}

	if listItemNode, ok := n.(*ast.ListItem); ok {
		// Increase indentation for list items
		level++
		for child := listItemNode.FirstChild(); child != nil; child = child.NextSibling() {
			processNode(child, reader, level, currentStyle, chunks)
		}
		return
	}

	// Process other children recursively
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		processNode(child, reader, level, currentStyle, chunks)
	}
}

// Main function to parse content
func parseContent(input string) []chunk {
	md := goldmark.New()
	reader := text.NewReader([]byte(input))
	document := md.Parser().Parse(reader)

	// Collect chunks
	var chunks []chunk
	processNode(document, reader, 0, encodeStyle(false, false), &chunks)
	return chunks
}

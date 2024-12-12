package slidesutils

import (
	"log"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Define the Chunk struct
type Chunk struct {
	Content          string
	Style            Style
	IndentationLevel int // 0 means no indentation, 1 first bullet, 2 second bullet
	Paragraph        int
}

// Define the custom type
type Style []byte

// Define bit masks for each style
const (
	NormalMask byte = 1 << 0 // 0000 0001
	BoldMask   byte = 1 << 1 // 0000 0010
	ItalicMask byte = 1 << 2 // 0000 0100
)

// Encode a style into a custom style type
func EncodeStyle(bold, italic bool) Style {
	var s byte
	if bold {
		s |= BoldMask
	}
	if italic {
		s |= ItalicMask
	}
	if !italic && !bold {
		s |= NormalMask
	}
	return Style{s}
}

// Decode a style back into its components
func DecodeStyle(s Style) (bold, italic, normal bool) {
	if len(s) == 0 {
		return false, false, false
	}
	bold = s[0]&BoldMask != 0
	italic = s[0]&ItalicMask != 0
	normal = s[0]&NormalMask != 0
	return
}

var paragraph int

// Recursive function to traverse AST and collect chunks
func processNode(n ast.Node, reader text.Reader, level int, currentStyle Style, chunks *[]Chunk) {
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
		*chunks = append(*chunks, Chunk{
			Content:          string(content),
			Style:            currentStyle,
			IndentationLevel: level,
			Paragraph:        paragraph,
		})
	}

	if emphasisNode, ok := n.(*ast.Emphasis); ok {
		// Update style for emphasis
		newStyle := currentStyle
		bold, italic, _ := DecodeStyle(currentStyle)
		if emphasisNode.Level == 1 {
			italic = true
		} else if emphasisNode.Level == 2 {
			bold = true
		}
		newStyle = EncodeStyle(bold, italic)

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
func ParseContent(input string) []Chunk {
	md := goldmark.New()
	reader := text.NewReader([]byte(input))
	document := md.Parser().Parse(reader)

	// Collect chunks
	var chunks []Chunk
	processNode(document, reader, 0, EncodeStyle(false, false), &chunks)
	return chunks
}

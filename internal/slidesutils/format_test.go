package slidesutils

import (
	"reflect"
	"testing"
)

func TestParseContent(t *testing.T) {
	input := `this is a **bold** word and this is an _italic_ like *this*. This is a list:
- the level of indentation should be 1
  - this content should have a level indentation of 2
and this is back to a level of indentation of zero`

	expected := []Chunk{
		{Content: "this is a ", Style: Style{NormalMask}, IndentationLevel: 0},
		{Content: "bold", Style: Style{BoldMask}, IndentationLevel: 0},
		{Content: " word and this is an ", Style: Style{NormalMask}, IndentationLevel: 0},
		{Content: "italic", Style: Style{ItalicMask}, IndentationLevel: 0},
		{Content: " like ", Style: Style{NormalMask}, IndentationLevel: 0},
		{Content: "this", Style: Style{ItalicMask}, IndentationLevel: 0},
		{Content: ". This is a list:", Style: Style{NormalMask}, IndentationLevel: 0},
		{Content: "the level of indentation should be 1\n", Style: Style{NormalMask}, IndentationLevel: 1},
		{Content: "this content should have a level indentation of 2\n", Style: Style{NormalMask}, IndentationLevel: 2},
		{Content: "and this is back to a level of indentation of zero", Style: Style{NormalMask}, IndentationLevel: 0},
	}

	result := parseContent(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("parseContent() = %+v, want %+v", result, expected)
	}
}

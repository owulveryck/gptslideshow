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
		{Content: "this is a ", Style: Style{NormalMask}, IndentationLevel: 0, Paragraph: 1},
		{Content: "bold", Style: Style{BoldMask}, IndentationLevel: 0, Paragraph: 1},
		{Content: " word and this is an ", Style: Style{NormalMask}, IndentationLevel: 0, Paragraph: 1},
		{Content: "italic", Style: Style{ItalicMask}, IndentationLevel: 0, Paragraph: 1},
		{Content: " like ", Style: Style{NormalMask}, IndentationLevel: 0, Paragraph: 1},
		{Content: "this", Style: Style{ItalicMask}, IndentationLevel: 0, Paragraph: 1},
		{Content: ". This is a list:", Style: Style{NormalMask}, IndentationLevel: 0, Paragraph: 1},
		{Content: "the level of indentation should be 1", Style: Style{NormalMask}, IndentationLevel: 1, Paragraph: 2},
		{Content: "this content should have a level indentation of 2", Style: Style{NormalMask}, IndentationLevel: 2, Paragraph: 3},
		{Content: "and this is back to a level of indentation of zero", Style: Style{NormalMask}, IndentationLevel: 0, Paragraph: 4},
	}

	result := parseContent(input)

	for i := 0; i < len(result); i++ {
		if !reflect.DeepEqual(expected[i], result[i]) {
			t.Errorf("want %+v, have %+v", expected[i], result[i])
		}
	}
}

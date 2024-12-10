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

	expected := []chunk{
		{content: "this is a ", style: style{normalMask}, indentationLevel: 0},
		{content: "bold", style: style{boldMask}, indentationLevel: 0},
		{content: " word and this is an ", style: style{normalMask}, indentationLevel: 0},
		{content: "italic", style: style{italicMask}, indentationLevel: 0},
		{content: " like ", style: style{normalMask}, indentationLevel: 0},
		{content: "this", style: style{italicMask}, indentationLevel: 0},
		{content: ". This is a list:", style: style{normalMask}, indentationLevel: 0},
		{content: "the level of indentation should be 1\n", style: style{normalMask}, indentationLevel: 1},
		{content: "this content should have a level indentation of 2\n", style: style{normalMask}, indentationLevel: 2},
		{content: "and this is back to a level of indentation of zero", style: style{normalMask}, indentationLevel: 0},
	}

	result := parseContent(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("parseContent() = %+v, want %+v", result, expected)
	}
}

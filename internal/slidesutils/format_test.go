package slidesutils

import (
	"reflect"
	"testing"
)

func TestParseContent(t *testing.T) {
	input := `this is a **bold** word and this is a list:
- the level of indentation should be 1
  - this content should have a level indentation of 2
and this is back to a level of indentation of zero`

	expected := []chunk{
		{content: "this is a ", isBold: false, indentationLevel: 0},
		{content: "bold", isBold: true, indentationLevel: 0},
		{content: " word and this is a list:\n", isBold: false, indentationLevel: 0},
		{content: "the level of indentation should be 1\n", isBold: false, indentationLevel: 1},
		{content: "this content should have a level indentation of 2\n", isBold: false, indentationLevel: 2},
		{content: "and this is back to a level of indentation of zero", isBold: false, indentationLevel: 0},
	}

	result := parseContent(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("parseContent() = %+v, want %+v", result, expected)
	}
}

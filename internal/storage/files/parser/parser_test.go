package parser_test

import (
	"testing"

	"github.com/tags-drive/core/internal/storage/files/parser"
)

func TestParse(t *testing.T) {
	tests := []struct {
		expr      string
		isCorrect bool
		answer    string
	}{
		// correct
		{"15", true, "15"},
		{"!6", true, "6 !"},
		{"7&8", true, "7 8 &"},
		{"1&2|3", true, "1 2 & 3 |"},
		{"!5&8|20", true, "5 ! 8 & 20 |"},
		{"66&!8|7", true, "66 8 ! & 7 |"},
		{"(1|2)&(2|3)", true, "1 2 | 2 3 | &"},
		{"(!7|6)&(6|9)", true, "7 ! 6 | 6 9 | &"},
		{"(8&9&10)|11", true, "8 9 & 10 & 11 |"},
		{"(!1|2)&2&3", true, "1 ! 2 | 2 & 3 &"},
		{"(!8|9)&6&66", true, "8 ! 9 | 6 & 66 &"},
		// incorrect
		{")5|8", false, ""},
		{"&", false, ""},
		{"|", false, ""},
		{"6&", false, ""},
		{"7&8|", false, ""},
		{"(!1|2))&2&3", false, ""},
		{"(!!1|3)&3&4", false, ""},
		{"(!&!2|3)&3&50", false, ""},
	}

	for i, tt := range tests {
		res, err := parser.Parse(tt.expr)
		if !tt.isCorrect && err == nil {
			t.Errorf("Test #%d Want: error Got: %s", i, res)
			continue
		}

		if tt.isCorrect && err != nil {
			t.Errorf("Test #%d Got error: %s", i, err)
			continue
		}

		if tt.answer != res {
			t.Errorf("Test #%d Want: %s Got: %s", i, tt.answer, res)
		}
	}
}

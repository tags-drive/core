package aggregation_test

import (
	"testing"

	"github.com/tags-drive/core/internal/storage/files/aggregation"
)

func TestIsGoodFile(t *testing.T) {
	tests := []struct {
		expr   string
		tags   []int
		answer bool
	}{
		{"", []int{}, true},
		{"", []int{16, 20}, true},
		{"15", []int{15}, true},
		{"15", []int{16}, false},
		{"6 !", []int{1, 2, 3, 7, 8, 9}, true},
		{"6 !", []int{1, 2, 3, 6, 7, 8, 9}, false},
		{"7 8 &", []int{7, 8}, true},
		{"7 8 &", []int{7, 5}, false},
		{"7 8 &", []int{6, 6}, false},
		{"66 8 ! & 7 |", []int{66, 5, 7}, true},    // 66&!8|7
		{"66 8 ! & 7 |", []int{66, 8, 5}, false},   // 66&!8|7
		{"66 8 ! & 7 |", []int{33, 1, 2}, false},   // 66&!8|7
		{"66 8 ! & 7 |", []int{66, 8, 7}, true},    // 66&!8|7
		{"7 ! 6 | 6 9 | &", []int{7, 6}, true},     // (!7|6)&(6|9)
		{"7 ! 6 | 6 9 | &", []int{7, 3, 9}, false}, // (!7|6)&(6|9)
		{"7 ! 6 | 6 9 | &", []int{7, 3, 2}, false}, // (!7|6)&(6|9)
		{"7 ! 6 | 6 9 | &", []int{6}, true},        // (!7|6)&(6|9)
	}

	for i, tt := range tests {
		res := aggregation.IsGoodFile(tt.expr, tt.tags)

		if res != tt.answer {
			t.Errorf("Test #%d Want: %t Got: %t", i, tt.answer, res)
		}
	}
}

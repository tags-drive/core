package storage_test

import (
	"testing"

	"github.com/ShoshinNikita/tags-drive/internal/storage"
)

func TestIsGoodFile(t *testing.T) {
	tests := []struct {
		m     storage.Mode
		fTags []string
		pTags []string
		res   bool
	}{
		{storage.AndMode, []string{"a", "b", "c"}, []string{"a", "c"}, true},
		{storage.AndMode, []string{"a", "b", "c"}, []string{"a", "e"}, false},
		{storage.OrMode, []string{"a", "b", "c"}, []string{"a", "e"}, true},
		{storage.OrMode, []string{"a", "b", "c"}, []string{"f", "e"}, false},
		{storage.NotMode, []string{"p", "b", "c"}, []string{"a", "e"}, true},
		{storage.NotMode, []string{"a", "b", "c"}, []string{"a", "e"}, false},
		// Empty file tags
		{storage.AndMode, []string{}, []string{"a", "e"}, false},
		{storage.OrMode, []string{}, []string{"a", "e"}, false},
		{storage.NotMode, []string{}, []string{"a", "e"}, true},
		// Empty passed tags
		{storage.AndMode, []string{"a", "b", "c"}, []string{}, true},
		{storage.OrMode, []string{"a", "b", "c"}, []string{}, true},
		{storage.NotMode, []string{"a", "b", "c"}, []string{}, true},
	}

	for i, tt := range tests {
		res := storage.IsGoodFile(tt.m, tt.fTags, tt.pTags)
		if res != tt.res {
			t.Errorf("Test #%d Want: %v Got %v", i, tt.res, res)
		}
	}

}

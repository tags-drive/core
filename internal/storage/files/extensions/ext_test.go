package extensions

import (
	"testing"
)

func TestGetExt(t *testing.T) {
	var (
		jpg, _ = allExtensions.get(".jpg")
		png, _ = allExtensions.get(".png")
	)

	tests := []struct {
		ext string
		res Ext
	}{
		{"jpg", jpg},
		{".jpg", jpg},
		{"JPG", jpg},
		{".JPG", jpg},
		{"png", png},
		{"PnG", png},
	}

	for i, tt := range tests {
		res := GetExt(tt.ext)
		if res != tt.res {
			t.Errorf("Test #%d Want: %v\nGot: %v\n", i, tt.res, res)
		}
	}
}

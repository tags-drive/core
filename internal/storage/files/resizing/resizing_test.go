package resizing_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/tags-drive/core/internal/storage/files/resizing"
)

func TestResizing(t *testing.T) {
	tests := []struct {
		origin  string
		resized string
	}{
		{"testdata/1.jpg", "testdata/1_res.jpg"},
		{"testdata/2.jpg", "testdata/2_res.jpg"},
		{"testdata/3.jpg", "testdata/3_res.jpg"},
		{"testdata/4.jpg", "testdata/4_res.jpg"},
		{"testdata/5.jpg", "testdata/5_res.jpg"},
		{"testdata/6.png", "testdata/6_res.png"},
	}
	for i, tt := range tests {
		file, err := os.Open(tt.origin)
		if err != nil {
			t.Errorf("Test #%d can't open file %s: %s", i, tt.origin, err)
		}
		im, err := resizing.Decode(file)
		file.Close()
		if err != nil {
			t.Fatal(err)
		}

		resizedImage := resizing.Resize(im)
		r, err := resizing.Encode(resizedImage, filepath.Ext(tt.origin))
		if err != nil {
			t.Fatal(err)
		}

		// Save an image
		f, err := os.Create(tt.resized)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(f, r)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()

		// Delete file
		err = os.Remove(tt.resized)
		if err != nil {
			t.Logf("Can't delete file %s: %s", tt.resized, err)
		}
	}
}

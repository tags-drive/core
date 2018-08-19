package web

import (
	"io"
	"net/http"
	"os"

	"github.com/ShoshinNikita/log"

	"github.com/ShoshinNikita/tags-drive/internal/params"
)

const (
	maxSize = 50000000 // 50MB
)

func upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		switch err {
		case http.ErrNotMultipart:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	for _, f := range r.MultipartForm.File["files"] {
		file, err := f.Open()
		if err != nil {
			log.Errorf("Can't open uploaded file: %s\n", err)
			continue
		}

		path := params.DataFolder + "/" + f.Filename
		// TODO
		if _, err := os.Open(path); !os.IsNotExist(err) {
			log.Errorf("File %s already exists\n", path)
			continue
		}

		newFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			log.Errorf("Can't create a new file %s: %s\n", path, err)
			continue
		}
		_, err = io.Copy(newFile, file)
		if err != nil {
			log.Errorf("Can't copy a new file %s: %s", path, err)
			// Deleting of a bad file
			err = os.Remove(path)
			if err != nil {
				log.Errorf("Can't delete a bad file %s: %s\n", path, err)
			}
			continue
		}
		log.Infof("File %s is uploaded\n", f.Filename)
	}
}

package web

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/ShoshinNikita/tags-drive/internal/storage"
)

const (
	maxSize = 50000000 // 50MB
)

// upload uploads files
//
// Request: multipart/form-data
// Response: json list of strings with status of files uploading
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

	msgChan := make(chan string, 50)
	wg := new(sync.WaitGroup)

	for _, f := range r.MultipartForm.File["files"] {
		wg.Add(1)
		go func(header *multipart.FileHeader) {
			defer wg.Done()
			err := storage.UploadFile(header, []string{})
			if err != nil {
				msgChan <- fmt.Sprintf("%s: %s", header.Filename, err)
			} else {
				msgChan <- fmt.Sprintf("%s: %s", header.Filename, "done")
			}
		}(f)
	}
	wg.Wait()
	close(msgChan)

	var messages []string

	for msg := range msgChan {
		messages = append(messages, msg)
	}

	json.NewEncoder(w).Encode(messages)
}

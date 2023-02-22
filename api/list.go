package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Error struct {
	Error string `json:"error"`
}

type File struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Created     time.Time `json:"created"`
	ContentType string    `json:"contentType"`
}

func getFileList() ([]File, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	q := &storage.Query{Prefix: ""}
	it := bucket.Objects(ctx, q)
	files := []File{}
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		file := File{Name: attrs.Name, Size: attrs.Size, Created: attrs.Created, ContentType: attrs.ContentType}
		files = append(files, file)
	}

	return files, nil
}

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	files, err := getFileList()
	if err != nil {
		e := Error{Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(e); err != nil {
			log.Fatal(err)
		}
	}

	if err := json.NewEncoder(w).Encode(files); err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}

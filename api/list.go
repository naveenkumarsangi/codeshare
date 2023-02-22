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
	Name   string    `json:"name"`
	Expiry time.Time `json:"expiry"`
}

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	q := &storage.Query{Prefix: ""}
	it := bucket.Objects(ctx, q)
	files := []File{}
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			e := Error{Error: err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(e); err != nil {
				log.Fatal(err)
			}
			return
		}

		file := File{Name: attrs.Name, Expiry: attrs.RetentionExpirationTime}
		files = append(files, file)
	}

	if err := json.NewEncoder(w).Encode(files); err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}

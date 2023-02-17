package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
)

// downloadFile downloads an object to a file.
func downloadFile(w http.ResponseWriter, bucket *storage.BucketHandle, object string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := bucket.Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	if _, err := io.Copy(w, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	return nil
}

func ShareHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	log.Println("id: ", id)

	if err := downloadFile(w, bucket, id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

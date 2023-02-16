package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var storageClient *storage.Client
var bucket *storage.BucketHandle

func getFirebaseApp() (*firebase.App, error) {
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), &firebase.Config{StorageBucket: "pastebin-personal.appspot.com"}, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	ctx := context.Background()
	storageClient, err = storage.NewClient(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("fbstorage.NewClient: %v", err)
	}

	return app, nil
}

func getDefaultBucket(app *firebase.App) (*storage.BucketHandle, error) {
	client, err := app.Storage(context.Background())
	if err != nil {
		return nil, err
	}

	return client.DefaultBucket()
}

// downloadFile downloads an object to a file.
func downloadFile(w http.ResponseWriter, bucket *storage.BucketHandle, object string) error {
	// object := "object-name"
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

// uploadFile uploads an object.
func uploadFile(data string, bucket *storage.BucketHandle, object string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := bucket.Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %v", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err := wc.Write([]byte(data)); err != nil {
		return fmt.Errorf("wc.Write: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	log.Printf("Blob %v uploaded.\n", object)
	return nil
}

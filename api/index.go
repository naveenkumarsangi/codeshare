package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func init() {
	fb, err := getFirebaseApp()
	if err != nil {
		log.Fatalf("getFirebaseApp: %v", err)
	}

	bucket, err = getDefaultBucket(fb)
	if err != nil {
		log.Fatalf("getDefaultBucket: %v", err)
	}
}

var bucket *storage.BucketHandle

func writeTempCredentials() (string, error) {
	sajson := os.Getenv("SERVICE_ACCOUNT_JSON")
	rawJson, err := base64.StdEncoding.DecodeString(sajson)
	if err != nil {
		return "", fmt.Errorf("error parsing service account json: %v", err)
	}

	f, err := os.CreateTemp("/tmp", "creds-*")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %v", err)
	}
	defer f.Close()

	if _, err := f.Write(rawJson); err != nil {
		return "", fmt.Errorf("error writing temp file: %v", err)
	}

	return f.Name(), nil
}

func getFirebaseApp() (*firebase.App, error) {
	credsFile, err := writeTempCredentials()
	if err != nil {
		log.Fatal(err)
	}

	opt := option.WithCredentialsFile(credsFile)
	app, err := firebase.NewApp(context.Background(), &firebase.Config{StorageBucket: "pastebin-personal.appspot.com"}, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
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

func PongHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

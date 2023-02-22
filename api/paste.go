package handler

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

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

func verifySignature(content, encsig string) error {
	sig, err := base64.StdEncoding.DecodeString(encsig)
	if err != nil {
		return err
	}

	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(content))
	if err != nil {
		return err
	}
	msgHashSum := msgHash.Sum(nil)

	return rsa.VerifyPKCS1v15(pubSigningKey, crypto.SHA256, msgHashSum, sig)
}

func PasteHandler(w http.ResponseWriter, r *http.Request) {
	sig := r.Header.Get("X-Api-Signature")
	if sig == "" {
		fmt.Fprintf(w, "missing signature")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		fmt.Fprintf(w, "error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	content := r.PostFormValue("content")
	if content == "" {
		fmt.Fprintf(w, "content in form is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := verifySignature(content, sig); err != nil {
		fmt.Fprintf(w, "crypto error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := uuid.New()
	if err := uploadFile(content, bucket, id.String()); err != nil {
		fmt.Fprintf(w, "unable to upload content, error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		scheme := "https"
		fmt.Fprintf(w, "%s://%s/api/share?id=%s", scheme, r.Host, id)
		w.WriteHeader(http.StatusCreated)
	}
}

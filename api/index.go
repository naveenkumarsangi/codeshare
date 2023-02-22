package handler

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const publicSigningKeyPEM = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAt1+KkWeWjBZ0djq6yMrc
mR5hQeWFwAl53tLtkWhWlO9PByCnWmmeKxAvSkVkjDzZwg+DOPkC1CXR2ZI1hGLk
/tgakXADN/U2u0pLYILDokE9hLQ+GQQP7CLohy7LQlw6uJRVGSdRwcuTrw3YIoIk
QXuzvaHH7+mZzWMq3bB8Pm1INBsRU8W/Hc4Y+tq95MJhiSiVi1TApjuX2L8ylF32
FHI3jjr1/Wkg75VHkj9mkw6wOsd9EgdteCzw9QrHX/21Fi6ODJ2KqnqGs9iLzQie
BDnBIXBVn7YU0P+i+LcBfLxT+pgSH7INcPgCl9UHLKxU3jCezwu/+XZ15tdC0xuH
w/AAA5frMcJZ2+eqXa0RuB4gwkoqnD+Oo6y2/KHjDQ6h65cSLtYhRmw7ewddUMXY
hNgIpYIQD7tmBGcX2guiwHdvIN+I3CftVFK3DzVakUxCOlDTQT1UQhs/MjYfQorh
9h+JHJX2HN1YqWm18wVneBl4YVrgGvjMQk+XDiDlhv0korB5SEhzEEgnsUwB2oyj
JBFkAr/Rr1CzVUVrhqE6VGTbHyKI83X7LiZQ+zV8GZ/C83kGkjoND6Qzc6rEUGn9
XBZVHWh8AF8HeWH+7L2AddjLaJqfbZKvmWpLkYjmSdfrgj04OQfY89SzvEMW5O3p
z1fWizapDUHSs3gQPvnv2h0CAwEAAQ==
-----END PUBLIC KEY-----`

func init() {
	fb, err := getFirebaseApp()
	if err != nil {
		log.Fatalf("getFirebaseApp: %v", err)
	}

	bucket, err = getDefaultBucket(fb)
	if err != nil {
		log.Fatalf("getDefaultBucket: %v", err)
	}

	block, _ := pem.Decode([]byte(publicSigningKeyPEM))
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("unable to decode PEM: %v", err)
	}

	if pubKey, ok := parsedKey.(*rsa.PublicKey); !ok {
		log.Fatalf("unable to parse RSA public key: %v", err)
	} else {
		pubSigningKey = pubKey
	}
}

var pubSigningKey *rsa.PublicKey
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

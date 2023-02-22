#!/usr/bin/env bash

function file_exists() {
	[ -f "$1" ] && echo "file already exists: $1" && exit 1
}

file_exists "$HOME/.ssh/codeshare.pem"
file_exists "$HOME/.ssh/codeshare.public.pem"

openssl genrsa -aes128 -passout pass:codeshare -out private.pem 4096
openssl rsa -in private.pem -passin pass:codeshare -pubout -out public.pem

mv private.pem ~/.ssh/codeshare.pem
mv public.pem ~/.ssh/codeshare.public.pem

echo "Please copy the following public key into api/index.go"
cat ~/.ssh/codeshare.public.pem

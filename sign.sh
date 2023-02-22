#!/usr/bin/env bash

TMPFILE="/tmp/codeshare.sign.sha256"
[ $# -ge 1 -a -f "$1" ] && input="$1" || input="-"

openssl dgst -sha256 -sign $HOME/.ssh/codeshare.pem -out "$TMPFILE" $input
openssl base64 -A -in "$TMPFILE" -out signature

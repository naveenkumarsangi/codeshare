#!/usr/bin/env bash

TMPFILE="/tmp/codeshare.sign.sha256"

function usage() {
	echo "usage: $0 [SIGNATURE] [ORIGINAL_FILE]"
	exit 1
}

if [ "$1" == "" ] || [ "$2" == "" ]; then
	usage
fi

openssl base64 -A -d -in "$1" -out /tmp/codeshare.sign.sha256
openssl dgst -sha256 -verify $HOME/.ssh/codeshare.public.pem -signature /tmp/codeshare.sign.sha256 "$2"

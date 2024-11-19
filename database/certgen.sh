#!/bin/bash

KEY_FILE="privkey.key"
CERT_FILE="fullchain.crt"

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout "$KEY_FILE" \
    -out "$CERT_FILE"

chmod 777 $KEY_FILE $CERT_FILE

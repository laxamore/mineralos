#!/bin/bash

if [ -z "$1" ]; then
    echo "how to use: ./generate-ca-cert.sh \"expired_days\""
else
    openssl req -x509 -newkey rsa:4096 -days $1 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/CN=\"\""

    echo "CA's self-signed certificate"
    openssl x509 -in ca-cert.pem -noout -text
fi
#!/bin/bash

if [ -z "$1" ]; then
    echo "how to use: ./generate-cert.sh \"expired_days\""
else
    openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out cert-req.pem  -subj "/CN=\"\""
    openssl x509 -req -in cert-req.pem -days $1 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out grpc-public-cert.pem -extfile cert.cnf

    echo "Server's signed certificate"
    openssl x509 -in grpc-public-cert.pem -noout -text
fi
#!/bin/sh

CGO_LDFLAGS='-lstdc++ -lzmq -lsodium ' CGO_ENABLED=1 GOOS=linux go build -o bin/client -a --ldflags '-extldflags "-static" -v' zmq/client/main.go
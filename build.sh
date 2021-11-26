#!/bin/sh

CGO_LDFLAGS='-lstdc++ -lzmq -lsodium ' CGO_ENABLED=1 GOOS=linux go build -o linux/mineralos/bin/mineralos-daemon -a --ldflags '-extldflags "-static" -v' daemon/mineralos-daemon/mineralos-daemon.go
CGO_LDFLAGS='-lstdc++ -lzmq -lsodium ' CGO_ENABLED=1 GOOS=linux go build -o linux/mineralos/bin/mineralos-daemon-stop -a --ldflags '-extldflags "-static" -v' daemon/mineralos-daemon-stop/mineralos-daemon-stop.go
CGO_LDFLAGS='-lstdc++ -lzmq -lsodium ' CGO_ENABLED=1 GOOS=linux go build -o bin/router -a --ldflags '-extldflags "-static" -v' router/router.go
CGO_LDFLAGS='-lstdc++ -lzmq -lsodium ' CGO_ENABLED=1 GOOS=linux go build -o bin/worker -a --ldflags '-extldflags "-static" -v' worker/worker.go
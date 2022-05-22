install:
	make grpc-gen-proto mineralos-daemon grpc-server mineralos-bin

grpc-gen-proto:
	protoc --proto_path=grpc --go_out=grpc --go_opt=paths=source_relative --go-grpc_out=grpc --go-grpc_opt=paths=source_relative mineralos_proto/mineralos.proto

mineralos-daemon:
	CGO_ENABLED=0 GOOS=linux go build -o linux/mineralos/bin/mineralos-daemon -a --ldflags '-extldflags "-static" -v' daemon/mineralos-daemon/logrotation.go daemon/mineralos-daemon/mineralos-daemon.go
	CGO_ENABLED=0 GOOS=linux go build -o linux/mineralos/bin/mineralos-daemon-stop -a --ldflags '-extldflags "-static" -v' daemon/mineralos-daemon-stop/mineralos-daemon-stop.go

grpc-server:
	CGO_ENABLED=0 GOOS=linux go build -o bin/server -a --ldflags '-extldflags "-static" -v' grpc/server.go

mineralos-bin:
	CGO_ENABLED=0 GOOS=linux go build -o linux/mineralos/bin/gpuraw -a --ldflags '-extldflags "-static" -v' linux/src/gpuraw/gpulist.go linux/src/gpuraw/gpuraw.go
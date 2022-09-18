all:
	make grpc-gen-proto mineralos-daemon mineralos-server gpuraw restapi

grpc-gen-proto:
	protoc --proto_path=config --go_out=config --go_opt=paths=source_relative --go-grpc_out=config --go-grpc_opt=paths=source_relative mineralos_proto/mineralos.proto

mineralos-daemon:
	CGO_ENABLED=0 GOOS=linux go build -o build/mineralos/bin/mineralos-daemon -a --ldflags '-extldflags "-static" -v' cmd/mineralos_daemon/*.go
	CGO_ENABLED=0 GOOS=linux go build -o build/mineralos/bin/mineralos-daemon-stop -a --ldflags '-extldflags "-static" -v' cmd/mineralos_daemon_stop/*.go

restapi:
	CGO_ENABLED=0 GOOS=linux go build -o build/mineralos/bin/restapi -a --ldflags '-extldflags "-static" -v' cmd/restapi/restapi.go

mineralos-server:
	CGO_ENABLED=0 GOOS=linux go build -o build/mineralos/bin/server -a --ldflags '-extldflags "-static" -v' cmd/mineralos_server/server.go

gpuraw:
	CGO_ENABLED=0 GOOS=linux go build -o build/mineralos/bin/gpuraw -a --ldflags '-extldflags "-static" -v' cmd/gpuraw/*.go
LOCAL_BIN:=$(CURDIR)/bin
BINARY_NAME=main

build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY_NAME} cmd/user/main.go

docker-build:
	docker buildx build --no-cache --platform linux/amd64 -t auth:1.0.0 .

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.1

install-protoc-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

lint:
	GOBIN=$(LOCAL_BIN) bin/golangci-lint run ./... --config .golangci.pipeline.yaml

get-protoc-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate-api:
	make generate-api-v1

generate-api-v1:
	mkdir -p pkg/user_v1
	protoc --proto_path api/user_v1 \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/user_v1/user.proto
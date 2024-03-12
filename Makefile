ifneq ($(ENV),)
	include $(ENV).env
endif

LOCAL_BIN:=$(CURDIR)/bin
BINARY_NAME=main
CONFIG=$(ENV).env
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(POSTGRES_PORT_LOCAL) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"
TLS_PATH=tls
TESTS_PATH=github.com/polshe-v/microservices_auth/internal/service/...,github.com/polshe-v/microservices_auth/internal/api/...
TESTS_ATTEMPTS=5
TESTS_COVERAGE_FILE=coverage.out

# #################### #
# DEPENDENCIES & TOOLS #
# #################### #

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.18.0
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@v3.3.6
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7

get-protoc-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

lint:
	GOBIN=$(LOCAL_BIN) $(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

generate-api:
	mkdir -p pkg/swagger
	make generate-api-v1
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'

generate-api-v1: check-env
	mkdir -p pkg/user_v1
	protoc --proto_path api/user_v1 --proto_path vendor.protogen \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/user_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
	--grpc-gateway_out=pkg/user_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
	api/user_v1/user.proto
	sed -i -e 's/{HTTP_HOST}/$(HTTP_HOST)/g' pkg/swagger/api.swagger.json
	sed -i -e 's/{HTTP_PORT}/$(HTTP_PORT)/g' pkg/swagger/api.swagger.json

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi

generate-mocks:
	go generate ./internal/repository
	go generate ./internal/service

generate-cert: $(TLS_PATH)/ca.key $(TLS_PATH)/ca.pem
	openssl genrsa -out $(TLS_PATH)/auth.key 4096
	openssl req -new -key $(TLS_PATH)/auth.key -config openssl.cnf -out $(TLS_PATH)/auth.csr
	openssl x509 -req -in $(TLS_PATH)/auth.csr -CA $(TLS_PATH)/ca.pem -CAkey $(TLS_PATH)/ca.key -extfile openssl.cnf -extensions req_ext -out $(TLS_PATH)/auth.pem -days 365 -sha256
	rm -rf $(TLS_PATH)/auth.csr

check-env:
ifeq ($(ENV),)
	$(error No environment specified)
endif

# ##### #
# TESTS #
# ##### #

test:
	go clean -testcache
	-go test ./... -v -covermode count -coverpkg=$(TESTS_PATH) -count $(TESTS_ATTEMPTS)

test-coverage:
	go clean -testcache
	-go test ./... -v -coverprofile=$(TESTS_COVERAGE_FILE).tmp -covermode count -coverpkg=$(TESTS_PATH) -count $(TESTS_ATTEMPTS)
	grep -v "mocks/" $(TESTS_COVERAGE_FILE).tmp > $(TESTS_COVERAGE_FILE)
	rm $(TESTS_COVERAGE_FILE).tmp
	go tool cover -html=$(TESTS_COVERAGE_FILE) -o coverage.html
	go tool cover -func=$(TESTS_COVERAGE_FILE) | grep "total"

# ##### #
# BUILD #
# ##### #

build-app:
	GOOS=linux GOARCH=amd64 go build -o $(LOCAL_BIN)/${BINARY_NAME} cmd/user/main.go

docker-build: docker-build-app docker-build-migrator

docker-build-app: check-env
	docker buildx build --no-cache --platform linux/amd64 -t auth:${APP_IMAGE_TAG} --build-arg="ENV=${ENV}" --build-arg="CONFIG=${CONFIG}" .

docker-build-migrator: check-env
	docker buildx build --no-cache --platform linux/amd64 -t migrator-auth:${MIGRATOR_IMAGE_TAG} -f migrator.Dockerfile --build-arg="ENV=${ENV}" .

# ###### #
# DEPLOY #
# ###### #

docker-deploy: check-env docker-build
	docker compose --env-file=$(ENV).env up -d

local-migration-status: check-env
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up: check-env
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down: check-env
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

# #### #
# STOP #
# #### #

docker-stop: check-env
	docker compose --env-file=$(ENV).env down

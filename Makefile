API_PATH 			= api/proto
PROTO_API_DIR 		= api/proto
PROTO_OUT_DIR 		= pkg/gen-tr-syst
PROTO_API_OUT_DIR	= ${PROTO_OUT_DIR}

.PHONY: gen/openapi
gen/openapi:
	oapi-codegen --old-config-style -generate types -o pkg/authorization-gateway-api/restapi_types.gen.go -package bffadmin pkg/authorization-gateway-api/restapi.swagger.yml

.PHONY: gen/proto
gen/proto:
	mkdir -p ${PROTO_OUT_DIR}
	protoc \
		-I ${API_PATH} \
		--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
        --go-grpc_out=$(PROTO_OUT_DIR)  \
        --go-grpc_opt=paths=source_relative \
		./${PROTO_API_DIR}/*.proto

.PHONY: build-reverse-proxy
build-reverse-proxy:
	GOOS=linux GOARCH=amd64 go build -o ${GOBIN}/${REVERSE_PROXY} ./cmd/${REVERSE_PROXY}/ || exit 1

.PHONY: docker-build
docker-build:
	docker build . -t ${IMG_TAG} -t ${LATEST_TAG}

.PHONY: run-migrate-local
run-migrate-local:
	sql-migrate up -env="local"

.PHONY: run-migrate-development
run-migrate-development:
	sql-migrate up -env="development"

.PHONY: lint
lint: proto/lint go/lint

.PHONY: go/lint
go/lint:
	golangci-lint run  --config=.golangci.yml --timeout=180s ./...

# use PROTOLINT_FLAGS=-fix for instant fixes
.PHONY: proto/lint
proto/lint:
	protolint lint $(PROTOLINT_FLAGS) $(PROTO_API_DIR)/*

.PHONY: go-mock-install
go-mock-install:
	GO111MODULE=on go get github.com/golang/mock/mockgen@v1.5.0

.PHONY: test
test:
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: check
check: gen tidy test lint

.PHONY: install-deps
install-deps:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

.PHONY: build
gen: build/project build/reverse-proxy

.PHONY: build/project
build/project:
	docker build .\
		-t $(IMAGE_PREFIC)$(PROJECT_NAME) \
		--build-arg NAILBAR_USER=$(NAILBAR_USER) \
		--build-arg NAILBAR_TOKEN=$(NAILBAR_TOKEN) \
		--build-arg PROJECT_NAME=$(PROJECT_NAME) \
		--build-arg GOPRIVATE=$(GOPRIVATE) \
		--build-arg PORT=$(PROJECT_PORT)

.PHONY: build/reverse-proxy
build/reverse-proxy:
	docker build .\
		-t $(IMAGE_PREFIC)$(REVERSE_PROXY) \
		--build-arg NAILBAR_USER=$(NAILBAR_USER) \
		--build-arg NAILBAR_TOKEN=$(NAILBAR_TOKEN) \
		--build-arg PROJECT_NAME=$(REVERSE_PROXY) \
		--build-arg GOPRIVATE=$(GOPRIVATE) \
		--build-arg PORT=$(REVERSE_PROXY_PORT)

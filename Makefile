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

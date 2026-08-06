.PHONY: proto

PROTO_DIR := ./proto
OUT_DIR   := ./proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

export PATH := $(PATH):$(shell go env GOPATH)/bin

proto: $(PROTO_FILES)
	@echo "Gerando arquivos Go a partir dos protos..."
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Concluído."


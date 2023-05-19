# ==============================================================================
# Makefile helper functions for tool
#

SHELL := /bin/bash
PROTOC := protoc

ifeq ($(origin PROTOC_GEN_PATH), undefined)
	PROTOC_GEN_PATH := $(ROOT_DIR)/tools/proto
endif

PROTOC_EXIST := $(shell type $(PROTOC) >/dev/null 2>&1 || { echo >&1 "not installed"; })

.PHONY: protoc.verify
protoc.verify:
ifneq (${PROTOC_EXIST},)
	$(error Please install protoc first)
endif
	@echo "===========> protoc verification passed"

.PHONY: proto.bin.verify
proto.bin.verify: protoc.verify
	$(if $(wildcard $(PROTOC_GEN_PATH)/bin/protoc-gen-gofast),,go build -o $(PROTOC_GEN_PATH)/bin/protoc-gen-gofast github.com/gogo/protobuf/protoc-gen-gofast)
	$(if $(wildcard $(PROTOC_GEN_PATH)/bin/protoc-gen-go),,go build -o $(PROTOC_GEN_PATH)/bin/protoc-gen-go github.com/golang/protobuf/protoc-gen-go)
	$(if $(wildcard $(PROTOC_GEN_PATH)/bin/protoc-gen-grpc-gateway),,go build -o $(PROTOC_GEN_PATH)/bin/protoc-gen-grpc-gateway github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)
	$(if $(wildcard $(PROTOC_GEN_PATH)/bin/protoc-gen-swagger),,go build -o $(PROTOC_GEN_PATH)/bin/protoc-gen-swagger github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger)
	$(if $(wildcard $(PROTOC_GEN_PATH)/bin/protoc-gen-discovery),,cd $(PROTOC_GEN_PATH)/protoc-gen-discovery && go build -o $(PROTOC_GEN_PATH)/bin/protoc-gen-discovery .)

.PHONY: proto.gen
proto.gen: proto.bin.verify 
	@echo "===========> Generate grpc source codes"
	@sh $(ROOT_DIR)/script/proto.sh

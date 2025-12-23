PROTO_DIR=proto
PB_DIR=client/proto

PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")

.PHONY: proto clean

proto:
	@echo "Generating protobuf files..."
	@mkdir -p ${PB_DIR}
	@protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Protobuf generation complete."

clean:
	@echo "Cleaning generated protobuf files..."
	@rm -rf $(PB_DIR)
	@echo "Clean complete."

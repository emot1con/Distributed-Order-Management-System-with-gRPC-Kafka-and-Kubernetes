# Makefile for gRPC Proto Generation
# 
# Usage:
#   make proto-all          # Generate all proto files
#   make proto-user         # Generate user service proto
#   make proto-product      # Generate product service proto
#   make proto-order        # Generate order service proto
#   make proto-payment      # Generate payment service proto
#   make proto-broker       # Generate broker service proto
#   make proto-install      # Install required tools
#   make proto-check        # Check prerequisites

.PHONY: proto-all proto-user proto-product proto-order proto-payment proto-broker proto-install proto-check help

# Default target
help:
	@echo "gRPC Proto Generator Makefile"
	@echo "============================"
	@echo ""
	@echo "Available targets:"
	@echo "  proto-all       Generate all proto files"
	@echo "  proto-user      Generate user service proto"
	@echo "  proto-product   Generate product service proto"
	@echo "  proto-order     Generate order service proto"
	@echo "  proto-payment   Generate payment service proto"
	@echo "  proto-broker    Generate broker service proto"
	@echo "  proto-install   Install required tools"
	@echo "  proto-check     Check prerequisites"
	@echo "  help           Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make proto-all"
	@echo "  make proto-user"
	@echo "  make proto-install"

# Check if required tools are installed
proto-check:
	@echo "Checking prerequisites..."
	@which protoc > /dev/null || (echo "✗ protoc not found. Install from https://github.com/protocolbuffers/protobuf/releases" && exit 1)
	@echo "✓ protoc found"
	@which protoc-gen-go > /dev/null || (echo "✗ protoc-gen-go not found. Run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest" && exit 1)
	@echo "✓ protoc-gen-go found"
	@which protoc-gen-go-grpc > /dev/null || (echo "✗ protoc-gen-go-grpc not found. Run: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest" && exit 1)
	@echo "✓ protoc-gen-go-grpc found"
	@echo "All prerequisites met!"

# Install required tools
proto-install:
	@echo "Installing gRPC tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✓ Tools installed successfully"

# Generate user service proto
proto-user: proto-check
	@echo "Generating user service proto files..."
	cd user/proto && protoc --go_out=. --go-grpc_out=. user.proto
	@echo "✓ User service proto files generated"

# Generate product service proto
proto-product: proto-check
	@echo "Generating product service proto files..."
	cd product/proto && protoc --go_out=. --go-grpc_out=. product.proto
	@echo "✓ Product service proto files generated"

# Generate order service proto
proto-order: proto-check
	@echo "Generating order service proto files..."
	cd order/proto && protoc --go_out=. --go-grpc_out=. order.proto
	@if [ -f order/proto/product.proto ]; then \
		cd order/proto && protoc --go_out=. --go-grpc_out=. product.proto; \
	fi
	@echo "✓ Order service proto files generated"

# Generate payment service proto
proto-payment: proto-check
	@echo "Generating payment service proto files..."
	cd payment/proto && protoc --go_out=. --go-grpc_out=. payment.proto
	@if [ -f payment/proto/order.proto ]; then \
		cd payment/proto && protoc --go_out=. --go-grpc_out=. order.proto; \
	fi
	@echo "✓ Payment service proto files generated"

# Generate broker service proto
proto-broker: proto-check
	@echo "Generating broker service proto files..."
	cd broker/proto && for proto in *.proto; do \
		echo "Generating $$proto..."; \
		protoc --go_out=. --go-grpc_out=. $$proto; \
	done
	@echo "✓ Broker service proto files generated"

# Generate all proto files
proto-all: proto-check
	@echo "Generating all proto files..."
	@echo ""
	@$(MAKE) proto-user
	@echo ""
	@$(MAKE) proto-product
	@echo ""
	@$(MAKE) proto-order
	@echo ""
	@$(MAKE) proto-payment
	@echo ""
	@$(MAKE) proto-broker
	@echo ""
	@echo "✓ All proto files generated successfully!"

# Clean generated files
proto-clean:
	@echo "Cleaning generated proto files..."
	find . -name "*.pb.go" -type f -delete
	find . -name "*_grpc.pb.go" -type f -delete
	@echo "✓ Generated files cleaned"

# Watch for changes and regenerate
proto-watch:
	@echo "Watching for proto file changes..."
	@echo "Press Ctrl+C to stop"
	@while true; do \
		inotifywait -r -e modify --include='.*\.proto$$' . 2>/dev/null && \
		echo "Proto file changed, regenerating..." && \
		$(MAKE) proto-all; \
	done || echo "inotifywait not found. Install inotify-tools for watch functionality."

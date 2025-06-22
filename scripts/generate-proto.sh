#!/bin/bash
# gRPC Proto Generator Script for Go
# Run this script to generate Go code from proto files

echo "gRPC Proto Generator for Go"
echo "=========================="

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc not found. Please install Protocol Compiler."
    echo "   Download from: https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "❌ protoc-gen-go not found."
    echo "   Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "❌ protoc-gen-go-grpc not found."
    echo "   Install with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

echo "✅ All prerequisites found!"
echo ""

# Function to generate proto files
generate_proto() {
    local service=$1
    local proto_dir=$2
    local proto_file=$3
    
    echo "🔄 Generating $service service..."
    
    if [ ! -f "$proto_dir/$proto_file" ]; then
        echo "❌ Proto file not found: $proto_dir/$proto_file"
        return 1
    fi
    
    cd "$proto_dir" || exit 1
    protoc --go_out=. --go-grpc_out=. "$proto_file"
    
    if [ $? -eq 0 ]; then
        echo "✅ $service service generated successfully"
    else
        echo "❌ Failed to generate $service service"
        return 1
    fi
    
    cd - > /dev/null
}

# Generate all services
echo "Starting proto generation..."
echo ""

# User service
generate_proto "User" "../user/proto" "user.proto"
echo ""

# Product service
generate_proto "Product" "../product/proto" "product.proto"
echo ""

# Order service
generate_proto "Order" "../order/proto" "order.proto"
# Check if product.proto exists in order service
if [ -f "../order/proto/product.proto" ]; then
    generate_proto "Order-Product" "../order/proto" "product.proto"
fi
echo ""

# Payment service
generate_proto "Payment" "../payment/proto" "payment.proto"
# Check if order.proto exists in payment service
if [ -f "../payment/proto/order.proto" ]; then
    generate_proto "Payment-Order" "../payment/proto" "order.proto"
fi
echo ""

# Broker service
echo "🔄 Generating Broker service..."
cd "../broker/proto" || exit 1
for proto in *.proto; do
    if [ -f "$proto" ]; then
        echo "  → Generating $proto..."
        protoc --go_out=. --go-grpc_out=. "$proto"
    fi
done
echo "✅ Broker service generated successfully"
cd - > /dev/null

echo ""
echo "🎉 All proto files generated successfully!"
echo ""
echo "Generated files:"
find .. -name "*.pb.go" -o -name "*_grpc.pb.go" | sort

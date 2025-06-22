# gRPC Proto Generator Script for Go (PowerShell)
# Run this script to generate Go code from proto files

Write-Host "gRPC Proto Generator for Go" -ForegroundColor Green
Write-Host "==========================" -ForegroundColor Green
Write-Host ""

# Check prerequisites
Write-Host "Checking prerequisites..." -ForegroundColor Yellow

try {
    $null = Get-Command protoc -ErrorAction Stop
    Write-Host "✅ protoc found" -ForegroundColor Green
} catch {
    Write-Host "❌ protoc not found. Please install Protocol Compiler." -ForegroundColor Red
    Write-Host "   Download from: https://github.com/protocolbuffers/protobuf/releases" -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit 1
}

$goPath = go env GOPATH
$protocGenGo = Join-Path $goPath "bin\protoc-gen-go.exe"
if (Test-Path $protocGenGo) {
    Write-Host "✅ protoc-gen-go found" -ForegroundColor Green
} else {
    Write-Host "❌ protoc-gen-go not found." -ForegroundColor Red
    Write-Host "   Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest" -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit 1
}

$protocGenGoGrpc = Join-Path $goPath "bin\protoc-gen-go-grpc.exe"
if (Test-Path $protocGenGoGrpc) {
    Write-Host "✅ protoc-gen-go-grpc found" -ForegroundColor Green
} else {
    Write-Host "❌ protoc-gen-go-grpc not found." -ForegroundColor Red
    Write-Host "   Install with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest" -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host ""
Write-Host "Starting proto generation..." -ForegroundColor Cyan
Write-Host ""

# Function to generate proto files
function Generate-Proto {
    param(
        [string]$ServiceName,
        [string]$ProtoDir,
        [string]$ProtoFile
    )
    
    Write-Host "🔄 Generating $ServiceName service..." -ForegroundColor Yellow
    
    $protoPath = Join-Path $ProtoDir $ProtoFile
    if (-not (Test-Path $protoPath)) {
        Write-Host "❌ Proto file not found: $protoPath" -ForegroundColor Red
        return $false
    }
    
    Push-Location $ProtoDir
    try {
        & protoc --go_out=. --go-grpc_out=. $ProtoFile
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ $ServiceName service generated successfully" -ForegroundColor Green
            return $true
        } else {
            Write-Host "❌ Failed to generate $ServiceName service" -ForegroundColor Red
            return $false
        }
    } finally {
        Pop-Location
    }
}

# Get script directory
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# User service
Generate-Proto "User" (Join-Path $scriptDir "..\user\proto") "user.proto"
Write-Host ""

# Product service
Generate-Proto "Product" (Join-Path $scriptDir "..\product\proto") "product.proto"
Write-Host ""

# Order service
Generate-Proto "Order" (Join-Path $scriptDir "..\order\proto") "order.proto"
$orderProductProto = Join-Path $scriptDir "..\order\proto\product.proto"
if (Test-Path $orderProductProto) {
    Generate-Proto "Order-Product" (Join-Path $scriptDir "..\order\proto") "product.proto"
}
Write-Host ""

# Payment service
Generate-Proto "Payment" (Join-Path $scriptDir "..\payment\proto") "payment.proto"
$paymentOrderProto = Join-Path $scriptDir "..\payment\proto\order.proto"
if (Test-Path $paymentOrderProto) {
    Generate-Proto "Payment-Order" (Join-Path $scriptDir "..\payment\proto") "order.proto"
}
Write-Host ""

# Broker service
Write-Host "🔄 Generating Broker service..." -ForegroundColor Yellow
$brokerProtoDir = Join-Path $scriptDir "..\broker\proto"
Push-Location $brokerProtoDir
try {
    $protoFiles = Get-ChildItem -Filter "*.proto"
    foreach ($protoFile in $protoFiles) {
        Write-Host "  → Generating $($protoFile.Name)..." -ForegroundColor Gray
        & protoc --go_out=. --go-grpc_out=. $protoFile.Name
    }
    Write-Host "✅ Broker service generated successfully" -ForegroundColor Green
} finally {
    Pop-Location
}

Write-Host ""
Write-Host "🎉 All proto files generated successfully!" -ForegroundColor Green
Write-Host ""

# List generated files
Write-Host "Generated files:" -ForegroundColor Cyan
$generatedFiles = Get-ChildItem -Path (Join-Path $scriptDir "..") -Recurse -Include "*.pb.go", "*_grpc.pb.go" | Sort-Object FullName
foreach ($file in $generatedFiles) {
    Write-Host "  → $($file.FullName)" -ForegroundColor Gray
}

Write-Host ""
Read-Host "Press Enter to exit"

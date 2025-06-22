@echo off
REM gRPC Proto Generator Script for Go (Windows)
REM Run this script to generate Go code from proto files

echo gRPC Proto Generator for Go
echo ==========================

REM Check if protoc is installed
protoc --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ protoc not found. Please install Protocol Compiler.
    echo    Download from: https://github.com/protocolbuffers/protobuf/releases
    pause
    exit /b 1
)

REM Check Go tools
for /f %%i in ('go env GOPATH') do set GOPATH=%%i

if not exist "%GOPATH%\bin\protoc-gen-go.exe" (
    echo ❌ protoc-gen-go not found.
    echo    Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    pause
    exit /b 1
)

if not exist "%GOPATH%\bin\protoc-gen-go-grpc.exe" (
    echo ❌ protoc-gen-go-grpc not found.
    echo    Install with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    pause
    exit /b 1
)

echo ✅ All prerequisites found!
echo.

echo Starting proto generation...
echo.

REM User service
echo 🔄 Generating User service...
cd /d "%~dp0..\user\proto"
if exist "user.proto" (
    protoc --go_out=. --go-grpc_out=. user.proto
    if %errorlevel% equ 0 (
        echo ✅ User service generated successfully
    ) else (
        echo ❌ Failed to generate User service
    )
) else (
    echo ❌ user.proto not found
)
echo.

REM Product service
echo 🔄 Generating Product service...
cd /d "%~dp0..\product\proto"
if exist "product.proto" (
    protoc --go_out=. --go-grpc_out=. product.proto
    if %errorlevel% equ 0 (
        echo ✅ Product service generated successfully
    ) else (
        echo ❌ Failed to generate Product service
    )
) else (
    echo ❌ product.proto not found
)
echo.

REM Order service
echo 🔄 Generating Order service...
cd /d "%~dp0..\order\proto"
if exist "order.proto" (
    protoc --go_out=. --go-grpc_out=. order.proto
    if %errorlevel% equ 0 (
        echo ✅ Order service generated successfully
    ) else (
        echo ❌ Failed to generate Order service
    )
) else (
    echo ❌ order.proto not found
)

REM Check for product.proto in order service
if exist "product.proto" (
    echo   → Generating product.proto for order service...
    protoc --go_out=. --go-grpc_out=. product.proto
)
echo.

REM Payment service
echo 🔄 Generating Payment service...
cd /d "%~dp0..\payment\proto"
if exist "payment.proto" (
    protoc --go_out=. --go-grpc_out=. payment.proto
    if %errorlevel% equ 0 (
        echo ✅ Payment service generated successfully
    ) else (
        echo ❌ Failed to generate Payment service
    )
) else (
    echo ❌ payment.proto not found
)

REM Check for order.proto in payment service
if exist "order.proto" (
    echo   → Generating order.proto for payment service...
    protoc --go_out=. --go-grpc_out=. order.proto
)
echo.

REM Broker service
echo 🔄 Generating Broker service...
cd /d "%~dp0..\broker\proto"
for %%f in (*.proto) do (
    echo   → Generating %%f...
    protoc --go_out=. --go-grpc_out=. %%f
)
echo ✅ Broker service generated successfully
echo.

echo 🎉 All proto files generated successfully!
echo.

REM Return to original directory
cd /d "%~dp0"

echo Generated files can be found in each service's proto directory.
pause

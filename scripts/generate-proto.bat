@echo off
setlocal enabledelayedexpansion

echo gRPC Proto Generator for Go
echo ============================

REM Check if protoc is available
protoc --version >nul 2>&1
if !errorlevel! neq 0 (
    echo Error: protoc not found. Please install Protocol Compiler.
    echo Download from: https://github.com/protocolbuffers/protobuf/releases
    pause
    exit /b 1
)

REM Check if Go tools are available
for /f %%i in ('go env GOPATH') do set GOPATH=%%i

if not exist "%GOPATH%\bin\protoc-gen-go.exe" (
    echo Error: protoc-gen-go not found
    echo Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    pause
    exit /b 1
)

if not exist "%GOPATH%\bin\protoc-gen-go-grpc.exe" (
    echo Error: protoc-gen-go-grpc not found
    echo Install with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    pause
    exit /b 1
)

echo All prerequisites found!
echo.

REM Parse command line arguments
set "SERVICE=%1"
set "ACTION=%1"

if "%ACTION%"=="all" goto :generate_all
if "%ACTION%"=="user" goto :generate_user
if "%ACTION%"=="product" goto :generate_product
if "%ACTION%"=="order" goto :generate_order
if "%ACTION%"=="payment" goto :generate_payment
if "%ACTION%"=="broker" goto :generate_broker
if "%ACTION%"=="help" goto :show_help
if "%ACTION%"=="" goto :show_help

:show_help
echo Usage:
echo   generate-proto.bat [service^|all^|help]
echo.
echo Services:
echo   user     - Generate user service proto files
echo   product  - Generate product service proto files
echo   order    - Generate order service proto files
echo   payment  - Generate payment service proto files
echo   broker   - Generate broker service proto files
echo   all      - Generate all proto files
echo   help     - Show this help
echo.
echo Examples:
echo   generate-proto.bat all
echo   generate-proto.bat user
echo   generate-proto.bat product
goto :end

:generate_user
echo Generating user service proto files...
cd /d "%~dp0..\user\proto"
protoc --go_out=. --go-grpc_out=. user.proto
if !errorlevel! equ 0 (
    echo ✓ User service proto files generated successfully
) else (
    echo ✗ Failed to generate user service proto files
)
goto :end

:generate_product
echo Generating product service proto files...
cd /d "%~dp0..\product\proto"
protoc --go_out=. --go-grpc_out=. product.proto
if !errorlevel! equ 0 (
    echo ✓ Product service proto files generated successfully
) else (
    echo ✗ Failed to generate product service proto files
)
goto :end

:generate_order
echo Generating order service proto files...
cd /d "%~dp0..\order\proto"
protoc --go_out=. --go-grpc_out=. order.proto
if exist "product.proto" (
    protoc --go_out=. --go-grpc_out=. product.proto
)
if !errorlevel! equ 0 (
    echo ✓ Order service proto files generated successfully
) else (
    echo ✗ Failed to generate order service proto files
)
goto :end

:generate_payment
echo Generating payment service proto files...
cd /d "%~dp0..\payment\proto"
protoc --go_out=. --go-grpc_out=. payment.proto
if exist "order.proto" (
    protoc --go_out=. --go-grpc_out=. order.proto
)
if !errorlevel! equ 0 (
    echo ✓ Payment service proto files generated successfully
) else (
    echo ✗ Failed to generate payment service proto files
)
goto :end

:generate_broker
echo Generating broker service proto files...
cd /d "%~dp0..\broker\proto"
for %%f in (*.proto) do (
    echo Generating %%f...
    protoc --go_out=. --go-grpc_out=. %%f
)
if !errorlevel! equ 0 (
    echo ✓ Broker service proto files generated successfully
) else (
    echo ✗ Failed to generate broker service proto files
)
goto :end

:generate_all
echo Generating all proto files...
echo.

call :generate_user
echo.
call :generate_product
echo.
call :generate_order
echo.
call :generate_payment
echo.
call :generate_broker
echo.
echo ✓ All proto files generated successfully!
goto :end

:end
echo.
echo Done!
pause

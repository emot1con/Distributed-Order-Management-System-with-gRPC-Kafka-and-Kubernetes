SHELL=cmd.exe
BROKER_BINARY=brokerApp
USER_BINARY=userApp
PRODUCT_BINARY=productApp
ORDER_BINARY=orderApp
PAYMENT_BINARY=paymentApp

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_user build_broker build_product build_order build_payment
	@echo Stopping docker images (if running...)
	docker-compose down
	@echo Building (when required) and starting docker images...
	docker-compose up --build -d
	@echo Docker images built and started!

## down: stop docker compose
down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo Building broker binary...
	chdir ..\broker && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${BROKER_BINARY} ./cmd/app
	@echo Done!

## build_user: builds the user binary as a linux executable
build_user:
	@echo Building user binary...
	chdir ..\user && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${USER_BINARY} ./cmd/app
	@echo Done!

## build_product: builds the product binary as a linux executable
build_product:
	@echo Building product binary...
	chdir ..\product && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${PRODUCT_BINARY} ./cmd/app
	@echo Done!

## build_order: builds the order binary as a linux executable
build_order:
	@echo Building order binary...
	chdir ..\order && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${ORDER_BINARY} ./cmd/app
	@echo Done!

## build_payment: builds the payment binary as a linux executable
build_payment:
	@echo Building payment binary...
	chdir ..\payment && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${PAYMENT_BINARY} ./cmd/app
	@echo Done!

## stop: stop the front end
stop:
	@echo Stopping front end...
	@taskkill /IM "${FRONT_END_BINARY}" /F
	@echo "Stopped front end!"
all: build

init:
	@echo "Initializing..."
	@$(MAKE) tool_download
	@$(MAKE) build

build:
	@echo "Building..."
	@go mod tidy
	@$(MAKE) tool_update
	@$(MAKE) proto_gen

proto_gen:
	@echo "Generating proto..."
	buf dep update && \
	buf lint && \
	buf generate

tool_update:
	@echo "Updating tools..."
	@go get -modfile=tools.mod -tool github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go get -modfile=tools.mod -tool github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go get -modfile=tools.mod -tool google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go get -modfile=tools.mod -tool google.golang.org/protobuf/cmd/protoc-gen-go@latest

tool_download:
	@echo "Downloading tools..."
	@go install -modfile=tools.mod tool
	@go install github.com/bufbuild/buf/cmd/buf@latest

# Module management commands
mod-tidy:
	@echo "Tidying module..."
	@go mod tidy

mod-vendor:
	@echo "Vendoring dependencies..."
	@go mod vendor

tag:
	@echo "Creating git tag..."
	@read -p "Enter version (e.g., v1.0.1): " version; \
	git tag $$version && \
	git push origin $$version

publish: build
	@echo "Publishing module..."
	@git add .
	@git commit -m "Update generated proto files" || true
	@git push origin main

# Integration with consumer services
update-consumers:
	@echo "Updating consumer services..."
	@if [ -d "../paymentsrv" ]; then \
		cd ../paymentsrv && go get github.com/escape-ship/protos@latest && go mod tidy; \
	fi
	@if [ -d "../accountsrv" ]; then \
		cd ../accountsrv && go get github.com/escape-ship/protos@latest && go mod tidy; \
	fi
	@if [ -d "../ordersrv" ]; then \
		cd ../ordersrv && go get github.com/escape-ship/protos@latest && go mod tidy; \
	fi
	@if [ -d "../productsrv" ]; then \
		cd ../productsrv && go get github.com/escape-ship/protos@latest && go mod tidy; \
	fi
	@if [ -d "../gatewaysrv" ]; then \
		cd ../gatewaysrv && go get github.com/escape-ship/protos@latest && go mod tidy; \
	fi

# Release workflow
release: build publish tag update-consumers
	@echo "Release complete!"

run:
	@echo "Running..."
	@./bin/$(shell basename $(PWD))

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

clean:
	rm -f bin/$(shell basename $(PWD))
	rm -rf gen/*

.PHONY: all init build proto_gen tool_update tool_download mod-tidy mod-vendor tag publish update-consumers release clean

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a shared Protocol Buffers definitions repository for the "escape-ship" e-commerce platform. It contains gRPC service definitions that are used across multiple microservices in the project.

## Architecture

The codebase follows a microservices architecture with the following services:

- **Account Service** (`account.proto`): Handles user authentication including Kakao OAuth integration, login, and registration
- **Product Service** (`product.proto`): Manages product catalog, categories, and product options
- **Order Service** (`order.proto`): Handles order creation and retrieval with order items
- **Payment Service** (`payment.proto`): Integrates with Kakao Pay for payment processing

Each service uses gRPC-Gateway for HTTP/JSON API generation alongside gRPC endpoints.

## Project Structure

```
protos/
├── account.proto    # Authentication and user management
├── order.proto      # Order management
├── payment.proto    # Payment processing (Kakao Pay)
├── product.proto    # Product catalog and options
└── gen/            # Generated code directory (empty in this repo)
```

## Service Dependencies

Each microservice (`accountsrv/`, `ordersrv/`, `paymentsrv/`, `productsrv/`, `gatewaysrv/`) in the parent directory has its own `proto/` directory that copies these definitions and generates Go code using buf.

## Common Development Tasks

### Protocol Buffer Generation

Each service directory has its own `buf.gen.yaml` configuration and uses buf for code generation. The typical workflow is:

1. Edit `.proto` files in this directory
2. Copy updated files to individual service `proto/` directories
3. Run `buf generate` in each service directory to regenerate Go code

### Service Patterns

- All services use `google.api.annotations.proto` for HTTP endpoint mapping
- Payment service additionally uses `protoc-gen-openapiv2` for OpenAPI documentation
- Go package paths follow the pattern: `github.com/escape-ship/proto/{service}/gen`

### API Conventions

- HTTP endpoints use RESTful patterns where possible
- POST endpoints typically use `body: "*"` for JSON payloads
- Korean comments are present in some message definitions
- Services use standard gRPC status codes and error handling

## Integration Notes

- The gateway service (`gatewaysrv/`) aggregates all service definitions
- Authentication tokens are handled through the Account service
- Order processing involves coordination between Product, Order, and Payment services
- Kakao integration is present in both Account (OAuth) and Payment (Pay) services
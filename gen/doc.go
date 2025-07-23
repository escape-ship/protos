// Package gen provides generated Go bindings for the Escape Ship e-commerce platform's gRPC services.
//
// # Overview
//
// This package contains Protocol Buffer generated code for a microservices-based e-commerce platform
// called "Escape Ship". It provides client and server stubs for four core services that handle
// authentication, product management, order processing, and payment integration.
//
// # Services
//
// The platform consists of four main gRPC services:
//
//   - AccountService: User authentication and Kakao OAuth integration
//   - ProductService: Product catalog management with categories and options
//   - OrderService: Order creation and retrieval with detailed item tracking
//   - PaymentService: Kakao Pay payment processing integration
//
// # Architecture
//
// All services are designed to work together in a microservices architecture:
//
//	┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
//	│  AccountService │    │ ProductService  │    │  OrderService   │
//	│                 │    │                 │    │                 │
//	│ • OAuth Login   │    │ • Categories    │    │ • Order Items   │
//	│ • Registration  │    │ • Product CRUD  │    │ • Order Status  │
//	│ • Token Auth    │    │ • Options       │    │ • Tracking      │
//	└─────────────────┘    └─────────────────┘    └─────────────────┘
//	         │                       │                       │
//	         └───────────────────────┼───────────────────────┘
//	                                 │
//	                    ┌─────────────────┐
//	                    │ PaymentService  │
//	                    │                 │
//	                    │ • Kakao Pay     │
//	                    │ • Ready/Approve │
//	                    │ • Cancel        │
//	                    └─────────────────┘
//
// # Authentication Flow
//
// The AccountService supports both traditional email/password authentication and Kakao OAuth:
//
//	// Traditional Login
//	client := NewAccountServiceClient(conn)
//	response, err := client.Login(ctx, &LoginRequest{
//	    Email:    "user@example.com",
//	    Password: "password123",
//	})
//
//	// Kakao OAuth Flow
//	urlResp, err := client.GetKakaoLoginURL(ctx, &GetKakaoLoginURLRequest{})
//	// User visits urlResp.LoginUrl
//	callbackResp, err := client.GetKakaoCallBack(ctx, &GetKakaoCallBackRequest{
//	    Code: "oauth_code_from_kakao",
//	})
//
// # Product Management
//
// Products are organized with categories and support configurable options:
//
//	// Get all products
//	products, err := productClient.GetProducts(ctx, &GetProductsRequest{})
//
//	// Get specific product with options
//	product, err := productClient.GetProductByID(ctx, &GetProductByIDRequest{
//	    Id: "product-123",
//	})
//	options, err := productClient.GetProductOptions(ctx, &GetProductOptionsRequest{
//	    Id: "product-123",
//	})
//
// # Order Processing
//
// Orders contain multiple items with product details and support various payment methods:
//
//	order := &InsertOrderRequest{
//	    UserId:          "user-123",
//	    OrderNumber:     "ORD-2024-001",
//	    Status:          "pending",
//	    TotalPrice:      50000,
//	    PaymentMethod:   "kakao_pay",
//	    ShippingAddress: "123 Main St, Seoul",
//	    Items: []*InsertOrderItem{
//	        {
//	            ProductId:      "prod-1",
//	            ProductName:    "Sample Product",
//	            ProductOptions: "Size: M, Color: Blue",
//	            ProductPrice:   25000,
//	            Quantity:       2,
//	        },
//	    },
//	}
//	result, err := orderClient.InsertOrder(ctx, order)
//
// # Payment Integration
//
// Kakao Pay integration follows the standard prepare-approve-complete flow:
//
//	// 1. Prepare payment
//	readyResp, err := paymentClient.KakaoReady(ctx, &KakaoReadyRequest{
//	    PartnerOrderId: "order-123",
//	    PartnerUserId:  "user-456",
//	    ItemName:       "Order Items",
//	    Quantity:       2,
//	    TotalAmount:    50000,
//	    TaxFreeAmount:  0,
//	})
//
//	// 2. User completes payment on Kakao
//	// 3. Approve payment
//	approveResp, err := paymentClient.KakaoApprove(ctx, &KakaoApproveRequest{
//	    Tid:            readyResp.Tid,
//	    PartnerOrderId: "order-123",
//	    PartnerUserId:  "user-456",
//	    PgToken:        "payment_token_from_kakao",
//	})
//
// # HTTP/JSON Gateway
//
// All services support both gRPC and HTTP/JSON through grpc-gateway annotations.
// The HTTP endpoints follow RESTful conventions:
//
//	Account Service:
//	  GET  /oauth/kakao/login     - Get Kakao login URL
//	  POST /oauth/kakao/callback  - Handle OAuth callback
//	  POST /login                 - Traditional login
//	  POST /register              - User registration
//
//	Product Service:
//	  GET  /products              - List all products
//	  GET  /products/{id}         - Get specific product
//	  POST /products              - Create new product
//	  POST /product/{id}/options  - Get product options
//
//	Order Service:
//	  POST /v1/order/insert       - Create new order
//	  GET  /v1/order              - Get all orders
//
//	Payment Service:
//	  POST /payment/kakao/ready   - Prepare Kakao payment
//	  POST /payment/kakao/approve - Approve Kakao payment
//	  POST /payment/kakao/cancel  - Cancel Kakao payment
//
// # Error Handling
//
// All services use standard gRPC status codes for error reporting. Common patterns include:
//
//   - InvalidArgument: Malformed requests or validation failures
//   - NotFound: Requested resources don't exist
//   - Unauthenticated: Authentication required or failed
//   - Internal: Server-side processing errors
//
// # Development
//
// This package is generated from Protocol Buffer definitions using buf and the standard
// Go gRPC toolchain. The source .proto files are maintained separately and should be
// used for any modifications to the service definitions.
//
// Generated code includes:
//   - Message types for all request/response structures
//   - Service client interfaces for calling remote services
//   - Service server interfaces for implementing services
//   - HTTP/JSON gateway reverse proxy code
//
// # Dependencies
//
// Key dependencies include:
//   - google.golang.org/grpc - gRPC runtime
//   - google.golang.org/protobuf - Protocol Buffer runtime
//   - github.com/grpc-ecosystem/grpc-gateway/v2 - HTTP/gRPC gateway
//   - google.golang.org/genproto - Google API annotations
//
// For complete API documentation and examples, see the individual service client
// interfaces and message type definitions in this package.
package gen
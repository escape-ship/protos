// Package gen contains generated Protocol Buffer and gRPC code for the Escape Ship e-commerce platform.
//
// # Overview
//
// This package provides Go bindings for all microservices in the Escape Ship platform:
//   - AccountService: User authentication and Kakao OAuth integration
//   - OrderService: Order management and tracking
//   - PaymentService: Kakao Pay integration and payment processing
//   - ProductService: Product catalog and options management
//
// # Services
//
// ## AccountService
//
// Handles user authentication, registration, and Kakao OAuth integration.
//
// Endpoints:
//   - GetKakaoLoginURL: Generate Kakao OAuth login URL
//   - GetKakaoCallBack: Handle Kakao OAuth callback
//   - Login: User authentication
//   - Register: User registration
//
// ## OrderService
//
// Manages order creation and retrieval.
//
// Endpoints:
//   - InsertOrder: Create a new order
//   - GetAllOrders: Retrieve all orders
//
// ## PaymentService
//
// Integrates with Kakao Pay for payment processing.
//
// Endpoints:
//   - KakaoReady: Prepare payment with Kakao Pay
//   - KakaoApprove: Approve payment transaction
//   - KakaoCancel: Cancel payment transaction
//
// ## ProductService
//
// Manages product catalog, categories, and product options.
//
// Endpoints:
//   - GetProducts: Retrieve all products
//   - GetProductByID: Get specific product details
//   - PostProducts: Create new product
//   - GetProductOptions: Get product options and values
//
// # Usage Examples
//
// ## Creating a gRPC Client
//
//	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer conn.Close()
//
//	client := gen.NewAccountServiceClient(conn)
//
// ## Making a gRPC Call
//
//	ctx := context.Background()
//	resp, err := client.GetKakaoLoginURL(ctx, &gen.GetKakaoLoginURLRequest{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Login URL:", resp.LoginUrl)
//
// ## Implementing a gRPC Server
//
//	type accountServer struct {
//		gen.UnimplementedAccountServiceServer
//	}
//
//	func (s *accountServer) GetKakaoLoginURL(
//		ctx context.Context,
//		req *gen.GetKakaoLoginURLRequest,
//	) (*gen.GetKakaoLoginURLResponse, error) {
//		return &gen.GetKakaoLoginURLResponse{
//			LoginUrl: "https://kauth.kakao.com/oauth/authorize?...",
//		}, nil
//	}
//
// # HTTP/JSON Gateway
//
// All services also support HTTP/JSON through gRPC-Gateway. The generated
// gateway code automatically handles the conversion between HTTP requests
// and gRPC calls.
//
// Example HTTP endpoints:
//   - GET /oauth/kakao/login
//   - POST /oauth/kakao/callback
//   - POST /login
//   - GET /products
//   - POST /payment/kakao/ready
//
// # Message Types
//
// The package includes all Protocol Buffer message types for request/response
// handling across services. All field names follow Go naming conventions
// (PascalCase for exported fields) while maintaining Protocol Buffer
// snake_case field definitions.
//
// # Error Handling
//
// All service methods return standard Go errors that can be checked using
// the status package from google.golang.org/grpc/status for gRPC-specific
// error codes and details.
//
// # Thread Safety
//
// All generated client types are safe for concurrent use by multiple
// goroutines. Server implementations should ensure their own thread safety.
//
// # Version Compatibility
//
// This package follows semantic versioning. Breaking changes to the API
// will result in a major version bump. Backward compatibility is maintained
// within major versions.
//
// Generated from Protocol Buffer definitions in the escape-ship/protos repository.
package gen
package gen

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// ExampleAccountServiceClient demonstrates how to create and use an AccountServiceClient.
//
// This example shows the complete flow of connecting to a gRPC server and making
// authentication-related calls.
func ExampleAccountServiceClient() {
	// Create a connection to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create the account service client
	client := NewAccountServiceClient(conn)
	ctx := context.Background()

	// Example 1: Get Kakao login URL
	loginResp, err := client.GetKakaoLoginURL(ctx, &GetKakaoLoginURLRequest{})
	if err != nil {
		log.Fatalf("GetKakaoLoginURL failed: %v", err)
	}
	fmt.Printf("Kakao Login URL: %s\n", loginResp.LoginUrl)

	// Example 2: Handle Kakao callback
	callbackResp, err := client.GetKakaoCallBack(ctx, &GetKakaoCallBackRequest{
		Code: "authorization_code_from_kakao",
	})
	if err != nil {
		log.Fatalf("GetKakaoCallBack failed: %v", err)
	}
	fmt.Printf("Access Token: %s\n", callbackResp.AccessToken)

	// Example 3: User login
	loginUserResp, err := client.Login(ctx, &LoginRequest{
		Email:    "user@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Printf("User Access Token: %s\n", loginUserResp.AccessToken)

	// Example 4: User registration
	registerResp, err := client.Register(ctx, &RegisterRequest{
		Email:    "newuser@example.com",
		Password: "newpassword123",
	})
	if err != nil {
		log.Fatalf("Register failed: %v", err)
	}
	fmt.Printf("Registration: %s\n", registerResp.Message)
}

// ExampleOrderServiceClient demonstrates order management operations.
func ExampleOrderServiceClient() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := NewOrderServiceClient(conn)
	ctx := context.Background()

	// Example 1: Create a new order
	orderResp, err := client.InsertOrder(ctx, &InsertOrderRequest{
		UserId:          "user123",
		OrderNumber:     "ORD-20240122-001",
		Status:          "PENDING",
		TotalPrice:      50000,
		Quantity:        2,
		PaymentMethod:   "KAKAO_PAY",
		ShippingFee:     3000,
		ShippingAddress: "서울시 강남구 테헤란로 123",
		PaidAt:          "2024-01-22T10:30:00Z",
		Items: []*InsertOrderItemRequest{
			{
				ProductId:    "prod001",
				ProductName:  "스마트폰 케이스",
				ProductPrice: 25000,
				Quantity:     2,
			},
		},
	})
	if err != nil {
		log.Fatalf("InsertOrder failed: %v", err)
	}
	fmt.Printf("Created order ID: %s\n", orderResp.Id)

	// Example 2: Get all orders
	ordersResp, err := client.GetAllOrders(ctx, &GetAllOrdersRequest{})
	if err != nil {
		log.Fatalf("GetAllOrders failed: %v", err)
	}
	fmt.Printf("Found %d orders\n", len(ordersResp.Orders))
	for _, order := range ordersResp.Orders {
		fmt.Printf("Order %s: %s (₩%d)\n", order.OrderNumber, order.Status, order.TotalPrice)
	}
}

// ExamplePaymentServiceClient demonstrates Kakao Pay integration.
func ExamplePaymentServiceClient() {
	conn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := NewPaymentServiceClient(conn)
	ctx := context.Background()

	// Example 1: Prepare Kakao Pay transaction
	readyResp, err := client.KakaoReady(ctx, &KakaoReadyRequest{
		PartnerOrderId: "ORDER123",
		PartnerUserId:  "USER123",
		ItemName:       "스마트폰 케이스",
		Quantity:       1,
		TotalAmount:    25000,
		TaxFreeAmount:  0,
	})
	if err != nil {
		log.Fatalf("KakaoReady failed: %v", err)
	}
	fmt.Printf("Payment ready. TID: %s\n", readyResp.Tid)
	fmt.Printf("Redirect URL: %s\n", readyResp.NextRedirectPcUrl)

	// Example 2: Approve payment
	approveResp, err := client.KakaoApprove(ctx, &KakaoApproveRequest{
		Tid:            readyResp.Tid,
		PartnerOrderId: "ORDER123",
		PartnerUserId:  "USER123",
		PgToken:        "payment_token_from_kakao",
	})
	if err != nil {
		log.Fatalf("KakaoApprove failed: %v", err)
	}
	fmt.Printf("Payment approved for order: %s\n", approveResp.PartnerOrderId)

	// Example 3: Cancel payment (if needed)
	cancelResp, err := client.KakaoCancel(ctx, &KakaoCancelRequest{
		PartnerOrderId:        "ORDER123",
		CancelAmount:          "25000",
		CancelTaxFreeAmount:   0,
		CancelVatAmount:       0,
		CancelAvailableAmount: 25000,
	})
	if err != nil {
		log.Fatalf("KakaoCancel failed: %v", err)
	}
	fmt.Printf("Payment cancelled for order: %s\n", cancelResp.PartnerOrderId)
}

// ExampleProductServiceClient demonstrates product catalog operations.
func ExampleProductServiceClient() {
	conn, err := grpc.Dial("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := NewProductServiceClient(conn)
	ctx := context.Background()

	// Example 1: Get all products
	productsResp, err := client.GetProducts(ctx, &GetProductsRequest{})
	if err != nil {
		log.Fatalf("GetProducts failed: %v", err)
	}
	fmt.Printf("Found %d products\n", len(productsResp.Products))

	// Example 2: Get specific product
	productResp, err := client.GetProductByID(ctx, &GetProductByIDRequest{
		Id: "prod001",
	})
	if err != nil {
		log.Fatalf("GetProductByID failed: %v", err)
	}
	product := productResp.Product
	fmt.Printf("Product: %s (₩%d)\n", product.Name, product.Price)

	// Example 3: Create new product
	createResp, err := client.PostProducts(ctx, &PostProductsRequest{
		Name:       "새로운 스마트폰 케이스",
		Categories: []string{"전자제품", "액세서리"},
		Price:      29000,
		ImageUrl:   "https://example.com/image.jpg",
	})
	if err != nil {
		log.Fatalf("PostProducts failed: %v", err)
	}
	fmt.Printf("Product created: %s\n", createResp.Message)

	// Example 4: Get product options
	optionsResp, err := client.GetProductOptions(ctx, &GetProductOptionsRequest{
		Id: "prod001",
	})
	if err != nil {
		log.Fatalf("GetProductOptions failed: %v", err)
	}
	fmt.Printf("Product %s has %d options\n", optionsResp.ProductId, len(optionsResp.Options))
	for _, option := range optionsResp.Options {
		fmt.Printf("Option: %s (%d values)\n", option.OptionName, len(option.Values))
	}
}

// ExampleGRPCServer demonstrates how to implement a gRPC server using the generated service interfaces.
func ExampleGRPCServer() {
	// Create a gRPC server
	server := grpc.NewServer()

	// Register services (you would implement these)
	// RegisterAccountServiceServer(server, &accountServiceImpl{})
	// RegisterOrderServiceServer(server, &orderServiceImpl{})
	// RegisterPaymentServiceServer(server, &paymentServiceImpl{})
	// RegisterProductServiceServer(server, &productServiceImpl{})

	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("gRPC server listening on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// ExampleHTTPGateway demonstrates how to set up an HTTP/JSON gateway
// that proxies requests to gRPC services.
func ExampleHTTPGateway() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux()

	// gRPC server endpoint
	endpoint := "localhost:50051"
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register service handlers (you would uncomment these with real implementations)
	// Example of how to register handlers:
	// err := RegisterAccountServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// if err != nil {
	//     log.Fatalf("Failed to register account service: %v", err)
	// }
	
	// err = RegisterOrderServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// if err != nil {
	//     log.Fatalf("Failed to register order service: %v", err)
	// }
	
	// err = RegisterPaymentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// if err != nil {
	//     log.Fatalf("Failed to register payment service: %v", err)
	// }
	
	// err = RegisterProductServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// if err != nil {
	//     log.Fatalf("Failed to register product service: %v", err)
	// }

	// Start HTTP server
	fmt.Println("HTTP Gateway server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to serve HTTP gateway: %v", err)
	}
}

// Example service implementation structure for AccountService
type accountServiceExample struct {
	UnimplementedAccountServiceServer
}

// GetKakaoLoginURL implements the AccountServiceServer interface.
// This is an example implementation showing the expected method signature.
func (s *accountServiceExample) GetKakaoLoginURL(
	ctx context.Context,
	req *GetKakaoLoginURLRequest,
) (*GetKakaoLoginURLResponse, error) {
	// Your implementation logic here
	return &GetKakaoLoginURLResponse{
		LoginUrl: "https://kauth.kakao.com/oauth/authorize?client_id=YOUR_APP_KEY&redirect_uri=YOUR_REDIRECT_URI&response_type=code",
	}, nil
}

// GetKakaoCallBack implements the AccountServiceServer interface.
// This is an example implementation showing OAuth callback handling.
func (s *accountServiceExample) GetKakaoCallBack(
	ctx context.Context,
	req *GetKakaoCallBackRequest,
) (*GetKakaoCallBackResponse, error) {
	// Your OAuth callback implementation logic here
	// 1. Exchange authorization code for access token
	// 2. Get user info from Kakao API
	// 3. Create or update user in your database
	// 4. Generate your own access/refresh tokens
	
	return &GetKakaoCallBackResponse{
		AccessToken:  "your_generated_access_token",
		RefreshToken: "your_generated_refresh_token",
		UserInfoJson: `{"id": 123, "name": "사용자", "email": "user@example.com"}`,
	}, nil
}

// Login implements the AccountServiceServer interface.
// This is an example implementation showing user authentication.
func (s *accountServiceExample) Login(
	ctx context.Context,
	req *LoginRequest,
) (*LoginResponse, error) {
	// Your login implementation logic here
	// 1. Validate email and password
	// 2. Check user credentials against database
	// 3. Generate access and refresh tokens
	
	return &LoginResponse{
		AccessToken:  "generated_access_token",
		RefreshToken: "generated_refresh_token",
	}, nil
}

// Register implements the AccountServiceServer interface.
// This is an example implementation showing user registration.
func (s *accountServiceExample) Register(
	ctx context.Context,
	req *RegisterRequest,
) (*RegisterResponse, error) {
	// Your registration implementation logic here
	// 1. Validate email format and password strength
	// 2. Check if user already exists
	// 3. Hash password and store user in database
	// 4. Send confirmation email (optional)
	
	return &RegisterResponse{
		Message: "Registration successful",
	}, nil
}
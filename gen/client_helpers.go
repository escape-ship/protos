package gen

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

// ClientConfig holds configuration for gRPC clients.
// This is a helper struct to simplify client creation with common options.
type ClientConfig struct {
	// Address is the server address (host:port)
	Address string
	
	// Insecure disables TLS when true
	Insecure bool
	
	// Timeout for connection establishment
	DialTimeout time.Duration
	
	// KeepAlive parameters
	KeepAliveTime    time.Duration
	KeepAliveTimeout time.Duration
	
	// Authority override for TLS
	Authority string
}

// DefaultClientConfig returns a configuration suitable for local development.
// For production, you should customize these settings.
func DefaultClientConfig(address string) *ClientConfig {
	return &ClientConfig{
		Address:          address,
		Insecure:         true, // Set to false in production
		DialTimeout:      5 * time.Second,
		KeepAliveTime:    30 * time.Second,
		KeepAliveTimeout: 5 * time.Second,
	}
}

// ProductionClientConfig returns a configuration suitable for production use
// with TLS enabled and appropriate timeouts.
func ProductionClientConfig(address string) *ClientConfig {
	return &ClientConfig{
		Address:          address,
		Insecure:         false,
		DialTimeout:      10 * time.Second,
		KeepAliveTime:    60 * time.Second,
		KeepAliveTimeout: 10 * time.Second,
	}
}

// NewConnection creates a new gRPC connection using the provided configuration.
// The caller is responsible for closing the connection.
//
// Example:
//
//	config := DefaultClientConfig("localhost:50051")
//	conn, err := NewConnection(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer conn.Close()
func NewConnection(config *ClientConfig) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	// Configure TLS
	if config.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		creds := credentials.NewTLS(&tls.Config{
			ServerName: config.Authority,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	// Configure keepalive
	keepaliveParams := keepalive.ClientParameters{
		Time:                config.KeepAliveTime,
		Timeout:             config.KeepAliveTimeout,
		PermitWithoutStream: true,
	}
	opts = append(opts, grpc.WithKeepaliveParams(keepaliveParams))

	// Connect with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()

	return grpc.DialContext(ctx, config.Address, opts...)
}

// AuthenticatedContext creates a context with authentication metadata.
// This is useful for adding JWT tokens or API keys to requests.
//
// Example:
//
//	ctx := AuthenticatedContext(context.Background(), "Bearer", "your-jwt-token")
//	resp, err := client.GetAllOrders(ctx, &GetAllOrdersRequest{})
func AuthenticatedContext(ctx context.Context, authType, token string) context.Context {
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %s", authType, token))
	return metadata.NewOutgoingContext(ctx, md)
}

// ClientSet provides convenient access to all service clients.
// This is useful when you need to make calls to multiple services.
type ClientSet struct {
	Account AccountServiceClient
	Order   OrderServiceClient
	Payment PaymentServiceClient
	Product ProductServiceClient
	conn    *grpc.ClientConn
}

// NewClientSet creates a new ClientSet using the provided configuration.
// All clients share the same underlying connection.
//
// Example:
//
//	config := DefaultClientConfig("localhost:50051")
//	clients, err := NewClientSet(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer clients.Close()
//
//	// Use any service client
//	resp, err := clients.Account.Login(ctx, &LoginRequest{...})
func NewClientSet(config *ClientConfig) (*ClientSet, error) {
	conn, err := NewConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	return &ClientSet{
		Account: NewAccountServiceClient(conn),
		Order:   NewOrderServiceClient(conn),
		Payment: NewPaymentServiceClient(conn),
		Product: NewProductServiceClient(conn),
		conn:    conn,
	}, nil
}

// Close closes the underlying gRPC connection.
// This should be called when the ClientSet is no longer needed.
func (cs *ClientSet) Close() error {
	return cs.conn.Close()
}

// ServiceAddresses holds the addresses for all services.
// This is useful when services are deployed on different hosts/ports.
type ServiceAddresses struct {
	Account string
	Order   string
	Payment string
	Product string
}

// DistributedClientSet provides clients for services running on different addresses.
// Unlike ClientSet, each service client has its own connection.
type DistributedClientSet struct {
	Account     AccountServiceClient
	Order       OrderServiceClient
	Payment     PaymentServiceClient
	Product     ProductServiceClient
	connections []*grpc.ClientConn
}

// NewDistributedClientSet creates clients for services running on different addresses.
// Each service gets its own connection, allowing for different configurations per service.
//
// Example:
//
//	addresses := &ServiceAddresses{
//		Account: "account-service:50051",
//		Order:   "order-service:50052",
//		Payment: "payment-service:50053",
//		Product: "product-service:50054",
//	}
//	
//	clients, err := NewDistributedClientSet(addresses, ProductionClientConfig)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer clients.Close()
func NewDistributedClientSet(addresses *ServiceAddresses, configFn func(string) *ClientConfig) (*DistributedClientSet, error) {
	var connections []*grpc.ClientConn
	
	// Create connection for account service
	accountConn, err := NewConnection(configFn(addresses.Account))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to account service: %w", err)
	}
	connections = append(connections, accountConn)

	// Create connection for order service
	orderConn, err := NewConnection(configFn(addresses.Order))
	if err != nil {
		// Close previously created connections on error
		accountConn.Close()
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}
	connections = append(connections, orderConn)

	// Create connection for payment service
	paymentConn, err := NewConnection(configFn(addresses.Payment))
	if err != nil {
		accountConn.Close()
		orderConn.Close()
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}
	connections = append(connections, paymentConn)

	// Create connection for product service
	productConn, err := NewConnection(configFn(addresses.Product))
	if err != nil {
		accountConn.Close()
		orderConn.Close()
		paymentConn.Close()
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}
	connections = append(connections, productConn)

	return &DistributedClientSet{
		Account:     NewAccountServiceClient(accountConn),
		Order:       NewOrderServiceClient(orderConn),
		Payment:     NewPaymentServiceClient(paymentConn),
		Product:     NewProductServiceClient(productConn),
		connections: connections,
	}, nil
}

// Close closes all underlying gRPC connections.
func (dcs *DistributedClientSet) Close() error {
	var lastErr error
	for _, conn := range dcs.connections {
		if err := conn.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// RetryConfig configures retry behavior for gRPC calls.
type RetryConfig struct {
	MaxAttempts int
	BackoffBase time.Duration
	BackoffMax  time.Duration
}

// DefaultRetryConfig returns a sensible default retry configuration.
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts: 3,
		BackoffBase: 100 * time.Millisecond,
		BackoffMax:  5 * time.Second,
	}
}

// Example usage patterns and common scenarios

// ExampleKakaoPaymentFlow demonstrates a complete Kakao Pay transaction flow.
// This example shows how to coordinate between different services for a payment.
func ExampleKakaoPaymentFlow() {
	// This is a comprehensive example showing service interaction
	config := DefaultClientConfig("localhost:50051")
	clients, err := NewClientSet(config)
	if err != nil {
		panic(err)
	}
	defer clients.Close()

	ctx := context.Background()
	
	// Step 1: User authentication (if needed)
	authCtx := AuthenticatedContext(ctx, "Bearer", "user-jwt-token")
	
	// Step 2: Get product information
	product, err := clients.Product.GetProductByID(authCtx, &GetProductByIDRequest{
		Id: "prod001",
	})
	if err != nil {
		panic(err)
	}
	
	// Step 3: Create order
	order, err := clients.Order.InsertOrder(authCtx, &InsertOrderRequest{
		UserId:        "user123",
		OrderNumber:   "ORD-20240122-001",
		Status:        "PENDING",
		TotalPrice:    product.Product.Price,
		Quantity:      1,
		PaymentMethod: "KAKAO_PAY",
		Items: []*InsertOrderItemRequest{
			{
				ProductId:    product.Product.Id,
				ProductName:  product.Product.Name,
				ProductPrice: product.Product.Price,
				Quantity:     1,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	
	// Step 4: Prepare Kakao Pay
	paymentReady, err := clients.Payment.KakaoReady(authCtx, &KakaoReadyRequest{
		PartnerOrderId: order.Id,
		PartnerUserId:  "user123",
		ItemName:       product.Product.Name,
		Quantity:       1,
		TotalAmount:    product.Product.Price,
		TaxFreeAmount:  0,
	})
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Payment ready. Redirect user to: %s\n", paymentReady.NextRedirectPcUrl)
	
	// Step 5: After user completes payment, approve it
	// (This would be called from your callback handler)
	paymentApprove, err := clients.Payment.KakaoApprove(authCtx, &KakaoApproveRequest{
		Tid:            paymentReady.Tid,
		PartnerOrderId: order.Id,
		PartnerUserId:  "user123",
		PgToken:        "payment_token_from_kakao_callback",
	})
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Payment completed for order: %s\n", paymentApprove.PartnerOrderId)
}
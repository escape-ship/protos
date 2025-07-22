package gen

// This file contains documentation for common patterns and best practices
// when using the escape-ship proto definitions.

// Common Usage Patterns
//
// # Error Handling Pattern
//
// Always check for gRPC-specific errors using the status package:
//
//	import "google.golang.org/grpc/status"
//	import "google.golang.org/grpc/codes"
//
//	resp, err := client.Login(ctx, req)
//	if err != nil {
//		if st, ok := status.FromError(err); ok {
//			switch st.Code() {
//			case codes.InvalidArgument:
//				// Handle invalid input
//			case codes.Unauthenticated:
//				// Handle authentication failure
//			case codes.Internal:
//				// Handle server error
//			default:
//				// Handle other errors
//			}
//		}
//		return err
//	}
//
// # Context Management Pattern
//
// Always use contexts with timeouts for gRPC calls:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	
//	resp, err := client.GetProducts(ctx, &GetProductsRequest{})
//
// # Authentication Pattern
//
// Use metadata for authentication across services:
//
//	import "google.golang.org/grpc/metadata"
//	
//	md := metadata.Pairs("authorization", "Bearer "+token)
//	ctx := metadata.NewOutgoingContext(context.Background(), md)
//	
//	resp, err := client.GetAllOrders(ctx, &GetAllOrdersRequest{})
//
// # Retry Pattern
//
// Implement exponential backoff for transient failures:
//
//	import "time"
//	import "google.golang.org/grpc/codes"
//	import "google.golang.org/grpc/status"
//	
//	func retryableCall[T any](fn func() (T, error)) (T, error) {
//		var result T
//		maxRetries := 3
//		baseDelay := 100 * time.Millisecond
//		
//		for i := 0; i < maxRetries; i++ {
//			result, err := fn()
//			if err == nil {
//				return result, nil
//			}
//			
//			if st, ok := status.FromError(err); ok {
//				switch st.Code() {
//				case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted:
//					// Retry these errors
//					if i < maxRetries-1 {
//						delay := time.Duration(1<<uint(i)) * baseDelay
//						time.Sleep(delay)
//						continue
//					}
//				default:
//					// Don't retry other errors
//					return result, err
//				}
//			}
//		}
//		return result, err
//	}
//
// # Validation Pattern
//
// Validate request data before making gRPC calls:
//
//	func validateLoginRequest(req *LoginRequest) error {
//		if req.Email == "" {
//			return status.Error(codes.InvalidArgument, "email is required")
//		}
//		if req.Password == "" {
//			return status.Error(codes.InvalidArgument, "password is required")
//		}
//		// Add email format validation
//		return nil
//	}
//
// # Connection Pool Pattern
//
// Use connection pooling for high-throughput applications:
//
//	import "sync"
//	
//	type ConnectionPool struct {
//		configs []ClientConfig
//		clients []*ClientSet
//		mutex   sync.RWMutex
//		current int
//	}
//	
//	func (p *ConnectionPool) GetClient() *ClientSet {
//		p.mutex.RLock()
//		defer p.mutex.RUnlock()
//		
//		client := p.clients[p.current]
//		p.current = (p.current + 1) % len(p.clients)
//		return client
//	}
//
// # Graceful Shutdown Pattern
//
// Always implement graceful shutdown for gRPC servers:
//
//	import "os"
//	import "os/signal"
//	import "syscall"
//	
//	func startServerWithGracefulShutdown(server *grpc.Server, lis net.Listener) {
//		// Start server in goroutine
//		go func() {
//			if err := server.Serve(lis); err != nil {
//				log.Printf("Server error: %v", err)
//			}
//		}()
//		
//		// Wait for interrupt signal
//		c := make(chan os.Signal, 1)
//		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//		<-c
//		
//		// Graceful shutdown
//		log.Println("Shutting down server...")
//		server.GracefulStop()
//	}
//
// # Health Check Pattern
//
// Implement health checks for service monitoring:
//
//	import "google.golang.org/grpc/health"
//	import "google.golang.org/grpc/health/grpc_health_v1"
//	
//	healthServer := health.NewServer()
//	grpc_health_v1.RegisterHealthServer(server, healthServer)
//	
//	// Set service status
//	healthServer.SetServingStatus("AccountService", grpc_health_v1.HealthCheckResponse_SERVING)
//
// # Load Balancing Pattern
//
// Use client-side load balancing with multiple endpoints:
//
//	import "google.golang.org/grpc/resolver"
//	import "google.golang.org/grpc/balancer/roundrobin"
//	
//	// Register static resolver for multiple endpoints
//	conn, err := grpc.Dial(
//		"static:///localhost:50051,localhost:50052,localhost:50053",
//		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//
// # Monitoring Pattern
//
// Add monitoring and metrics to your gRPC services:
//
//	import "google.golang.org/grpc/stats"
//	
//	type statsHandler struct{}
//	
//	func (h *statsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
//		// Add tags for metrics
//		return ctx
//	}
//	
//	func (h *statsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
//		// Record metrics
//	}
//	
//	// Add to server
//	server := grpc.NewServer(grpc.StatsHandler(&statsHandler{}))

// Best Practices for Protocol Buffer Messages
//
// # Field Validation
//
// Always validate required fields in your service implementations:
//
//	func (s *server) CreateOrder(ctx context.Context, req *InsertOrderRequest) (*InsertOrderResponse, error) {
//		if req.UserId == "" {
//			return nil, status.Error(codes.InvalidArgument, "user_id is required")
//		}
//		if req.TotalPrice <= 0 {
//			return nil, status.Error(codes.InvalidArgument, "total_price must be positive")
//		}
//		// ... implementation
//	}
//
// # Timestamp Handling
//
// Use proper timestamp formats and timezone handling:
//
//	import "time"
//	
//	// When setting timestamps in requests
//	now := time.Now().UTC().Format(time.RFC3339)
//	req.PaidAt = now
//	
//	// When parsing timestamps from responses
//	paidAt, err := time.Parse(time.RFC3339, resp.PaidAt)
//	if err != nil {
//		return fmt.Errorf("invalid timestamp format: %w", err)
//	}
//
// # Money Handling
//
// Always use integer types for money to avoid floating-point precision issues:
//
//	// Store amounts in smallest currency unit (cents for USD, won for KRW)
//	// ₩1,000 = 1000 (in won)
//	// $10.50 = 1050 (in cents)
//	
//	func formatPrice(priceInWon int64) string {
//		return fmt.Sprintf("₩%s", humanize.Comma(priceInWon))
//	}
//
// # Pagination Pattern
//
// Implement pagination for list operations:
//
//	// Add to your proto definitions:
//	// message GetProductsRequest {
//	//   int32 page_size = 1;
//	//   string page_token = 2;
//	// }
//	// 
//	// message GetProductsResponse {
//	//   repeated Product products = 1;
//	//   string next_page_token = 2;
//	// }
//	
//	func (s *server) GetProducts(ctx context.Context, req *GetProductsRequest) (*GetProductsResponse, error) {
//		pageSize := req.PageSize
//		if pageSize <= 0 || pageSize > 100 {
//			pageSize = 20 // default page size
//		}
//		
//		// Implement pagination logic
//		products, nextToken, err := s.productRepo.GetPaginated(req.PageToken, pageSize)
//		if err != nil {
//			return nil, status.Error(codes.Internal, "failed to get products")
//		}
//		
//		return &GetProductsResponse{
//			Products:      products,
//			NextPageToken: nextToken,
//		}, nil
//	}

// Performance Optimization Patterns
//
// # Connection Reuse
//
// Reuse gRPC connections instead of creating new ones for each request:
//
//	// Good: Reuse connection
//	var globalClient AccountServiceClient
//	
//	func init() {
//		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
//		if err != nil {
//			panic(err)
//		}
//		globalClient = NewAccountServiceClient(conn)
//	}
//	
//	// Bad: Create new connection each time
//	// func makeCall() {
//	//     conn, _ := grpc.Dial(...)
//	//     client := NewAccountServiceClient(conn)
//	//     // ...
//	//     conn.Close()
//	// }
//
// # Streaming Pattern
//
// Use streaming for large data transfers:
//
//	// For bidirectional streaming (if implemented in your services)
//	stream, err := client.StreamOrders(ctx)
//	if err != nil {
//		return err
//	}
//	
//	// Send requests
//	go func() {
//		for _, req := range requests {
//			stream.Send(req)
//		}
//		stream.CloseSend()
//	}()
//	
//	// Receive responses
//	for {
//		resp, err := stream.Recv()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			return err
//		}
//		// Process response
//	}
//
// # Batch Operations Pattern
//
// Batch multiple operations to reduce network overhead:
//
//	// Instead of making multiple individual calls
//	// for _, productId := range productIds {
//	//     product, err := client.GetProductByID(ctx, &GetProductByIDRequest{Id: productId})
//	//     // ...
//	// }
//	
//	// Use batch operations when available
//	// batch, err := client.GetProductsBatch(ctx, &GetProductsBatchRequest{Ids: productIds})

// Security Patterns
//
// # Input Sanitization
//
// Sanitize all user input before processing:
//
//	import "html"
//	import "strings"
//	
//	func sanitizeInput(input string) string {
//		// Remove dangerous characters
//		input = strings.TrimSpace(input)
//		input = html.EscapeString(input)
//		return input
//	}
//	
//	func (s *server) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
//		email := sanitizeInput(req.Email)
//		if !isValidEmail(email) {
//			return nil, status.Error(codes.InvalidArgument, "invalid email format")
//		}
//		// ... implementation
//	}
//
// # Rate Limiting Pattern
//
// Implement rate limiting to prevent abuse:
//
//	import "golang.org/x/time/rate"
//	
//	type rateLimitedServer struct {
//		limiter *rate.Limiter
//		AccountServiceServer
//	}
//	
//	func (s *rateLimitedServer) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
//		if !s.limiter.Allow() {
//			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
//		}
//		return s.AccountServiceServer.Login(ctx, req)
//	}
//
// # Audit Logging Pattern
//
// Log all important operations for security auditing:
//
//	func (s *server) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
//		// Log login attempt
//		log.Printf("Login attempt for email: %s from IP: %s", 
//			req.Email, getClientIP(ctx))
//		
//		resp, err := s.authenticateUser(req.Email, req.Password)
//		if err != nil {
//			log.Printf("Login failed for email: %s - %v", req.Email, err)
//			return nil, err
//		}
//		
//		log.Printf("Login successful for email: %s", req.Email)
//		return resp, nil
//	}
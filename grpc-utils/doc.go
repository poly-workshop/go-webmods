// Package grpc_utils provides gRPC server interceptors for logging and request
// ID tracking with context propagation.
//
// # Interceptors
//
// This package provides two main interceptors:
//   - BuildLogInterceptor: Structured logging for gRPC requests
//   - BuildRequestIDInterceptor: Request ID generation and propagation
//
// # Log Interceptor
//
// The log interceptor provides structured logging for all gRPC requests:
//
//	import (
//	    "log/slog"
//	    "github.com/poly-workshop/go-webmods/grpc_utils"
//	    "google.golang.org/grpc"
//	)
//
//	logger := slog.Default()
//	server := grpc.NewServer(
//	    grpc.ChainUnaryInterceptor(
//	        grpc_utils.BuildLogInterceptor(logger),
//	    ),
//	)
//
// The interceptor logs:
//   - Request method
//   - Request duration
//   - Response status
//   - Any errors
//
// # Request ID Interceptor
//
// The request ID interceptor ensures every request has a unique ID:
//
//	server := grpc.NewServer(
//	    grpc.ChainUnaryInterceptor(
//	        grpc_utils.BuildRequestIDInterceptor(),
//	    ),
//	)
//
// Request ID handling:
//   - Checks incoming metadata for existing x-request-id
//   - Generates a new UUID if not present
//   - Adds request_id to log context via app.WithLogAttrs
//   - Returns x-request-id in response headers
//
// This enables request tracing across service boundaries.
//
// # Combined Usage
//
// Typically, both interceptors are used together:
//
//	import (
//	    "log/slog"
//	    "github.com/poly-workshop/go-webmods/app"
//	    "github.com/poly-workshop/go-webmods/grpc_utils"
//	    "google.golang.org/grpc"
//	)
//
//	func main() {
//	    app.SetCMDName("myservice")
//	    app.Init(".")
//
//	    logger := slog.Default()
//	    server := grpc.NewServer(
//	        grpc.ChainUnaryInterceptor(
//	            grpc_utils.BuildRequestIDInterceptor(),
//	            grpc_utils.BuildLogInterceptor(logger),
//	        ),
//	    )
//
//	    // Register your services
//	    // pb.RegisterYourServiceServer(server, &yourService{})
//
//	    // Start server...
//	}
//
// Order matters: Request ID interceptor should run first to ensure the ID
// is available for logging.
//
// # Request Tracing
//
// The request ID propagates through the system:
//
// 1. Client sends request (optionally with x-request-id header)
// 2. Server interceptor extracts or generates request ID
// 3. Server adds request ID to context
// 4. All logs within the request include the request ID
// 5. Server returns request ID in response headers
//
// Client-side request ID propagation:
//
//	import (
//	    "google.golang.org/grpc/metadata"
//	)
//
//	// Add request ID to outgoing context
//	md := metadata.Pairs("x-request-id", requestID)
//	ctx := metadata.NewOutgoingContext(ctx, md)
//
//	// Make gRPC call
//	resp, err := client.SomeMethod(ctx, req)
//
//	// Extract request ID from response
//	var header metadata.MD
//	resp, err = client.SomeMethod(ctx, req, grpc.Header(&header))
//	if ids := header.Get("x-request-id"); len(ids) > 0 {
//	    requestID := ids[0]
//	}
//
// # Integration with app Package
//
// The interceptors integrate with the app package for enhanced logging:
//
//	// In your gRPC handler
//	func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
//	    // Add additional context
//	    ctx = app.WithLogAttrs(ctx, slog.String("user_id", req.UserId))
//
//	    slog.InfoContext(ctx, "Fetching user")
//	    // Logs include: cmd, hostname, request_id, user_id
//
//	    // ... fetch user logic
//	}
//
// # Structured Logging
//
// All logs use structured key-value pairs:
//
//	slog.InfoContext(ctx, "Request processed",
//	    "method", "GetUser",
//	    "user_id", userID,
//	    "duration_ms", duration.Milliseconds(),
//	)
//
// Automatic fields in every log:
//   - cmd: Service name (from app.SetCMDName)
//   - hostname: Current hostname
//   - request_id: Unique request identifier
//   - Any fields added via app.WithLogAttrs
//
// # Best Practices
//
//   - Always use both interceptors together
//   - Place BuildRequestIDInterceptor before BuildLogInterceptor
//   - Propagate request IDs across service calls
//   - Use structured logging with key-value pairs
//   - Add relevant context via app.WithLogAttrs
//   - Log at appropriate levels (Info for normal, Error for failures)
//   - Include timing information in logs for performance monitoring
//
// # Error Handling
//
// The interceptors handle errors gracefully:
//   - Failed header setting logs an error but doesn't fail the request
//   - Missing request ID generates a new one
//   - Logging errors don't interrupt request processing
//
// # Performance Considerations
//
//   - UUID generation is fast but not free - use existing IDs when available
//   - Structured logging has minimal overhead with slog
//   - Interceptors add negligible latency to requests
//
// # Example Server Setup
//
// Complete example with health checks and reflection:
//
//	import (
//	    "net"
//	    "log/slog"
//	    "github.com/poly-workshop/go-webmods/app"
//	    "github.com/poly-workshop/go-webmods/grpc_utils"
//	    "google.golang.org/grpc"
//	    "google.golang.org/grpc/health"
//	    "google.golang.org/grpc/health/grpc_health_v1"
//	    "google.golang.org/grpc/reflection"
//	)
//
//	func main() {
//	    app.SetCMDName("api-server")
//	    app.Init(".")
//
//	    lis, err := net.Listen("tcp", ":50051")
//	    if err != nil {
//	        panic(err)
//	    }
//
//	    logger := slog.Default()
//	    server := grpc.NewServer(
//	        grpc.ChainUnaryInterceptor(
//	            grpc_utils.BuildRequestIDInterceptor(),
//	            grpc_utils.BuildLogInterceptor(logger),
//	        ),
//	    )
//
//	    // Register services
//	    // pb.RegisterYourServiceServer(server, &yourService{})
//
//	    // Health check
//	    healthServer := health.NewServer()
//	    grpc_health_v1.RegisterHealthServer(server, healthServer)
//	    healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
//
//	    // Reflection for grpcurl
//	    reflection.Register(server)
//
//	    slog.Info("Server starting", "address", ":50051")
//	    if err := server.Serve(lis); err != nil {
//	        panic(err)
//	    }
//	}
package grpc_utils

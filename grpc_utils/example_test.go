package grpc_utils_test

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/poly-workshop/go-webmods/app"
	"github.com/poly-workshop/go-webmods/grpc_utils"
	"google.golang.org/grpc"
)

// Example demonstrates setting up a gRPC server with logging and request ID interceptors.
func Example() {
	// Initialize application
	app.SetCMDName("grpc-server")
	// app.Init(".")

	// Create logger
	logger := slog.Default()

	// Create gRPC server with interceptors
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_utils.BuildRequestIDInterceptor(),
			grpc_utils.BuildLogInterceptor(logger),
		),
	)

	// Register your services
	// pb.RegisterYourServiceServer(server, &yourService{})

	// Start server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	fmt.Println("Server configured with interceptors")
	_ = server
	_ = lis
	// server.Serve(lis)

	// Output: Server configured with interceptors
}

// Example_requestIDPropagation demonstrates how request IDs propagate through the system.
func Example_requestIDPropagation() {
	// In your gRPC handler
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// The request ID is automatically added to the context by the interceptor
		// Add additional context
		ctx = app.WithLogAttrs(ctx, slog.String("user_id", "123"))

		// All logs will include request_id, user_id, cmd, and hostname
		slog.InfoContext(ctx, "Processing request")

		// The request ID is automatically returned in response headers
		return nil, nil
	}

	_ = handler
	fmt.Println("Request ID automatically propagated")
	// Output: Request ID automatically propagated
}

// Example_contextLogging demonstrates adding context to logs within gRPC handlers.
func Example_contextLogging() {
	// In your gRPC handler
	processUser := func(ctx context.Context, userID string) {
		// Add user context
		ctx = app.WithLogAttrs(ctx,
			slog.String("user_id", userID),
			slog.String("operation", "update"),
		)

		slog.InfoContext(ctx, "Starting user update")

		// All subsequent logs with this context include user_id and operation
		slog.InfoContext(ctx, "Validating user data")
		slog.InfoContext(ctx, "Updating database")
		slog.InfoContext(ctx, "User update complete")
	}

	_ = processUser
	fmt.Println("Context-aware logging configured")
	// Output: Context-aware logging configured
}

// Example_fullServer demonstrates a complete gRPC server setup.
func Example_fullServer() {
	// Initialize application
	// app.SetCMDName("api-server")
	// app.Init(".")

	// Create server with all interceptors
	// logger := slog.Default()
	// server := grpc.NewServer(
	//     grpc.ChainUnaryInterceptor(
	//         grpc_utils.BuildRequestIDInterceptor(),
	//         grpc_utils.BuildLogInterceptor(logger),
	//     ),
	// )

	// Register services
	// pb.RegisterYourServiceServer(server, &yourService{})

	// Add health check
	// healthServer := health.NewServer()
	// grpc_health_v1.RegisterHealthServer(server, healthServer)
	// healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Add reflection for grpcurl
	// reflection.Register(server)

	// Start server
	// lis, err := net.Listen("tcp", ":50051")
	// if err != nil {
	//     panic(err)
	// }

	// slog.Info("Server starting", "address", ":50051")
	// if err := server.Serve(lis); err != nil {
	//     panic(err)
	// }

	fmt.Println("Full server setup example")
	// Output: Full server setup example
}

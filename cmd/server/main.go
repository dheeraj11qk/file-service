package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"file-service/internal/db"
	"file-service/internal/interceptor"
	"file-service/internal/observability"
	"file-service/internal/service"
	"file-service/proto"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {

	// 1️⃣ Initialize Tracer FIRST
	shutdown := observability.InitTracer("file-service")
	defer shutdown(context.Background())
	// 1️⃣ Initialize Database
	db.InitDB()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Metrics server running on :9091")
		http.ListenAndServe(":9091", nil)
	}()

	// 2️⃣ Create TCP listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	secret := os.Getenv("JWT_SECRET")

	// 3️⃣ Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(
			interceptor.AuthInterceptor(secret),
		),
	)

	// 4️⃣ Register service
	proto.RegisterFileServiceServer(grpcServer, &service.FileService{})

	log.Println("🚀 gRPC Server running on port 50051")

	// 5️⃣ Start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}

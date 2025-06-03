package main

import (
	"log"
	"net"
	"net/http"

	pb "github.com/AnitsuriW/ollama-healthcheck/proto"
	"github.com/AnitsuriW/ollama-healthcheck/server"

	"google.golang.org/grpc"
)

func main() {
	// 启动 RESTful 服务
	go func() {
		http.HandleFunc("/health", server.HealthHandler)
		log.Println("REST server listening on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// 启动 gRPC 服务
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterHealthServiceServer(grpcServer, &server.HealthServer{})

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

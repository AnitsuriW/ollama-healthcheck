package server

import (
	"context"

	pb "github.com/AnitsuriW/ollama-healthcheck/proto"
)

type HealthServer struct {
	pb.UnimplementedHealthServiceServer
}

func (s *HealthServer) Check(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	healthy, msg := CheckOllamaHealth()
	return &pb.HealthResponse{
		Healthy: healthy,
		Message: msg,
	}, nil
}

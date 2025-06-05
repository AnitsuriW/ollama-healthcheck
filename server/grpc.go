package server

import (
	"context"
	"time"

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

func (s *HealthServer) PredictFailure(ctx context.Context, req *pb.PredictRequest) (*pb.PredictResponse, error) {
	// 动态预测逻辑
	var prediction string
	var confidence float64

	if req.CpuUsage > 80 || req.MemoryUsage > 85 {
		prediction = "failure_likely"
		confidence = 0.9
	} else if req.ResponseLatencyMs > 1000 || req.ErrorsLastMinute > 1 {
		prediction = "failure_possible"
		confidence = 0.7
	} else {
		prediction = "failure_unlikely"
		confidence = 0.3
	}

	return &pb.PredictResponse{
		Prediction: prediction,
		Confidence: float32(confidence),
		Timestamp:  time.Now().Format(time.RFC3339),
	}, nil
}

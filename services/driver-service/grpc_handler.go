package main

import (
	"context"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedDriverServiceServer
	service *Service
}

func NewgRPCHandler(server *grpc.Server, service *Service) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}
	pb.RegisterDriverServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method RegisterDriver not implemented")
}

func (h *gRPCHandler) UnRegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method UnRegisterDriver not implemented")
}

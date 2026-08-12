package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

// creating a constructior
func NewgRPCHandler(server *grpc.Server, service domain.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}
	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()

	pickUpCoordinates := &types.Coordinate{
		Latitude:  pickup.Latitude,
		Longitude: pickup.Longitude,
	}

	destinationCoordinates := &types.Coordinate{
		Latitude:  destination.Latitude,
		Longitude: destination.Longitude,
	}

	// get userID
	userID := req.GetUserID()

	route, err := h.service.GetRoute(ctx, pickUpCoordinates, destinationCoordinates)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	// 1. Estimate the ride fares prices based on the route (example distance)
	estimatedFares := h.service.EstimatePackagesPriceWithRoute(route)
	// 2. Store the ride Fares for the createtrip
	fares, err := h.service.GenerateTripFares(ctx, estimatedFares, userID, route)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate price: %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: domain.ToRideFareProto(fares),
	}, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	// task to do here
	fareID := req.GetRideFareID()
	userID := req.GetUserID()
	// 1) Fetch and validate the fare
	rideFare, err := h.service.GetAndValidate(ctx, fareID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create a trip %v", err)
	}
	// 2) Call create trip
	t, err := h.service.CreateTrip(ctx, rideFare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create trip: %v", err)
	}
	// 3) Also initialize and empty driver to the trip
	// done this in create trip struct

	// 4) add a comment at the end of the function to publish an event on the Async Comms module
	return &pb.CreateTripResponse{
		TripID: t.ID.Hex(),
	}, nil
}

package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var GrpcAddr = ":9093"

func main() {
	inmemRepo := repository.NewinmemRepository()
	svc := service.NewService(inmemRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// go routine to listen for signal ctrl+c or docker compose down etc
	go func() {
		// channel that receives os signals
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-sigCh
		cancel()
	}()

	// creating a grpc server
	// 1) Create a tcp connection
	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// 2) create grpc server
	grpcServer := grpcserver.NewServer()
	// TODO initialize our grpc handler implementation
	grpc.NewgRPCHandler(grpcServer, svc)

	log.Printf("Starting gRPC server, Trip service on port %s", lis.Addr().String())

	// 3) Start the server on goroutine
	// go routine to keep the server running
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Printf("shutting down the server: %v", err)
	grpcServer.GracefulStop()
}

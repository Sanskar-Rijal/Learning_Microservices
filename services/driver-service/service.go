package main

import (
	pb "ride-sharing/shared/proto/driver"
)

type Service struct {
	driver []*driverInMap
}

type driverInMap struct {
	Driver *pb.Driver
	// INDEX int
	// TODO route
}

func NewService() *Service {
	return &Service{
		driver: make([]*driverInMap, 0),
	}
}

// TODO register and unregister method

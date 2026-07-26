package repository

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
)

type inmemoryRepository struct {
	trips map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel 
}

func NewinmemRepository() *inmemoryRepository{
	return &inmemoryRepository{
		trips: map[string]*domain.TripModel{},
		rideFares: map[string]*domain.RideFareModel{},
	}
}

func(r *inmemoryRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error){

	r.trips[trip.ID.Hex()] = trip;
	return trip, nil; 
}
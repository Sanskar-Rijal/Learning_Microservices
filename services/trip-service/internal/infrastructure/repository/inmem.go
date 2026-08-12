package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
	"sync"
)

type inmemoryRepository struct {
	// creating mutex to lock and unlock while writing in memory
	mu        sync.RWMutex
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewinmemRepository() *inmemoryRepository {
	return &inmemoryRepository{
		trips:     map[string]*domain.TripModel{},
		rideFares: map[string]*domain.RideFareModel{},
	}
}

func (r *inmemoryRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (r *inmemoryRepository) SaveRideFare(ctx context.Context, fare *domain.RideFareModel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rideFares[fare.ID.Hex()] = fare
	return nil
}

func (r *inmemoryRepository) GetRideFareByID(ctx context.Context, fareID string) (*domain.RideFareModel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// find ridefare by id
	riderfare, exists := r.rideFares[fareID]

	if !exists {
		return nil, fmt.Errorf("fare doesn't exist with id: %s", fareID)
	}
	return riderfare, nil
}

package repository

import "ride-sharing/tools/services/trip-service/internal/domain"

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


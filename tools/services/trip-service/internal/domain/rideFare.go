package domain

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct{
	ID primitive.ObjectID
	UserID string
	PackageSlug string
	TotalPrice float64
	ExpiresAt time.Time
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error);
}
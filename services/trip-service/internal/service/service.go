package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}

	return s.repo.CreateTrip(ctx, t)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmApiResponse, error) {
	baseUrl := "http://router.project-osrm.org"

	url := fmt.Sprintf("%s/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		baseUrl,
		pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from osrm api: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}

	var routeResponse tripTypes.OsrmApiResponse

	if err := json.Unmarshal(body, &routeResponse); err != nil {
		return nil, fmt.Errorf("failed to parse the data: %v", err)
	}
	return &routeResponse, nil
}

func (s *service) EstimatePackagesPriceWithRoute(route *tripTypes.OsrmApiResponse) []*domain.RideFareModel {
	baseFares := getBaseFaresWithPackages()
	estimatedFares := make([]*domain.RideFareModel, len(baseFares))

	for i, f := range baseFares {
		estimatedFares[i] = estimateFareRoute(f, route)
	}

	return estimatedFares
}

func (s *service) GenerateTripFares(ctx context.Context, rideFares []*domain.RideFareModel, userID string) ([]*domain.RideFareModel, error) {
	// function to store in database
	fares := make([]*domain.RideFareModel, len(rideFares))

	for i, f := range rideFares {
		id := primitive.NewObjectID()
		fare := &domain.RideFareModel{
			UserID:          userID,
			ID:              id,
			TotalPriceInNpr: f.TotalPriceInNpr,
			PackageSlug:     f.PackageSlug,
		}

		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("Failed to save trip fare %s", &err)
		}
		fares[i] = fare
	}

	return fares, nil
}

// algorithm to calculate price for the ride
func estimateFareRoute(f *domain.RideFareModel, route *tripTypes.OsrmApiResponse) *domain.RideFareModel {
	// we will be considering
	// 1)Distance, 2) Time, 3) Car Price
	pricingCfg := tripTypes.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInNpr

	distanceKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	// distance fare
	distanceFare := distanceKm * pricingCfg.PricePerUnitOfDistance
	// fime fare
	timeFare := durationInMinutes * pricingCfg.PricingPerMinute

	// total price
	totalPriceinNPR := carPackagePrice + distanceFare + timeFare

	return &domain.RideFareModel{
		TotalPriceInNpr: totalPriceinNPR,
		PackageSlug:     f.PackageSlug,
	}
}

func getBaseFaresWithPackages() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug:     "suv",
			TotalPriceInNpr: 100,
		},
		{
			PackageSlug:     "sedan",
			TotalPriceInNpr: 200,
		},
		{
			PackageSlug:     "van",
			TotalPriceInNpr: 400,
		},
		{
			PackageSlug:     "luxury",
			TotalPriceInNpr: 800,
		},
	}
}

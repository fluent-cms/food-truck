package services

import (
	"context"
	"food-trucks/packages/models"
)

type FacilityStore interface {
	Set(ctx context.Context, vals []models.Facility) error
	Get(ctx context.Context, keys []string) ([]models.Facility, error)
}

type ItemFacilityStore interface {
	AddMem(ctx context.Context, sliceID any, items []models.Facility) error
	GetAllMemberEntities(ctx context.Context, sliceID any) ([]models.Facility, error)
}

type GeoFacilityStore interface {
	Add(ctx context.Context, item models.Facility) error
	Get(ctx context.Context, lat float64, lon float64, radius float64) ([]models.Facility, error)
}

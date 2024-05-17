package services

import (
	"context"
	"fmt"
	"food-trucks/packages/models"
	"food-trucks/packages/util/rdb"
	"testing"
)

func TestFacilitySvc_GetByLocation(t *testing.T) {
	ctx := context.Background()
	svc := mustInit()
	items, err := svc.GetByLocation(ctx, 37.805885350100986, -122.41594524663745, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(items)
}

func mustInit() *FacilitySvc {
	config := rdb.Config{}
	facilityStore := rdb.NewEntityStore[string, models.Facility]("facility", 0, config).
		WithGetKey(models.GetFacilityKey)
	itemFacilityStore := rdb.NewSliceStore[string, models.Facility]("item", 0, config, facilityStore).
		WithGetKey(models.GetFacilityKey).WithGetScore(models.GetFacilityScore)
	geoFacilityStore := rdb.NewGeoStore[string, models.Facility]("geo", 0, config, facilityStore).
		WithGetKey(models.GetFacilityKey).WithGetLocation(models.GetFacilityLocation)
	facilitySvc := &FacilitySvc{
		FacilityStore:     facilityStore,
		ItemFacilityStore: itemFacilityStore,
		GeoFacilityStore:  geoFacilityStore,
	}
	err := facilitySvc.Seed("../..//configs/data.csv")
	if err != nil {
		panic(err)
	}
	return facilitySvc
}

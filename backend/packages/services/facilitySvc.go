package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"food-trucks/packages/models"
	"food-trucks/packages/util/annotate"
	"os"
	"strconv"
	"strings"
)

type Location struct {
	Lat float64
	Lon float64
}
type FacilitySvc struct {
	FacilityStore     FacilityStore
	ItemFacilityStore ItemFacilityStore
	GeoFacilityStore  GeoFacilityStore
	Center            Location
}

func (t *FacilitySvc) GetByID(ctx context.Context, id string) (models.Facility, error) {
	var ret models.Facility
	items, err := t.FacilityStore.Get(ctx, []string{id})
	if err != nil {
		return ret, err
	}
	if len(items) == 0 {
		return ret, fmt.Errorf("not found facility " + id)
	}
	return items[0], nil
}

func (t *FacilitySvc) GetByItem(ctx context.Context, item string) ([]models.Facility, error) {
	return t.ItemFacilityStore.GetAllMemberEntities(ctx, strings.TrimSpace(item))
}

func (t *FacilitySvc) GetByLocation(ctx context.Context, lat, lon, radius float64) ([]models.Facility, error) {
	if lon == 0 || lat == 0 && radius == 0 {
		lat = t.Center.Lat
		lon = t.Center.Lon
		radius = 1
	}
	return t.GeoFacilityStore.Get(ctx, lat, lon, radius)
}

func (t *FacilitySvc) Seed(p string) error {
	ctx := context.Background()
	facilities, err := t.readCSV(p)
	if err != nil {
		return err
	}
	t.getCenter(facilities)
	if err = t.cacheFacilities(ctx, facilities); err != nil {
		return annotate.Errorf("failed to cache facilities, %w", err)
	}
	if err = t.cacheFoodItems(ctx, facilities); err != nil {
		return annotate.Errorf("failed to cache food items, %w", err)
	}
	return t.cacheLocations(ctx, facilities)
}

func (t *FacilitySvc) getCenter(facilities []models.Facility) {
	var lon, lat float64
	for _, facility := range facilities {
		lon += facility.Longitude
		lat += facility.Latitude
	}
	t.Center.Lat = lat / float64(len(facilities))
	t.Center.Lon = lon / float64(len(facilities))
}

func (t *FacilitySvc) cacheLocations(ctx context.Context, facilities []models.Facility) error {
	for _, facility := range facilities {
		if err := t.GeoFacilityStore.Add(ctx, facility); err != nil {
			return err
		}
	}
	return nil
}

func (t *FacilitySvc) cacheFoodItems(ctx context.Context, facilities []models.Facility) error {
	for _, facility := range facilities {
		for _, item := range strings.Split(facility.FoodItems, ":") {
			item = strings.TrimSpace(item)
			if err := t.ItemFacilityStore.AddMem(ctx, item, []models.Facility{facility}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *FacilitySvc) cacheFacilities(ctx context.Context, facilities []models.Facility) error {
	return t.FacilityStore.Set(ctx, facilities)
}

func (t *FacilitySvc) readCSV(p string) ([]models.Facility, error) {
	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // allows variable number of fields per record
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var facilities []models.Facility

	// Skip the header row
	for i, record := range records {
		if i == 0 {
			continue
		}

		x, _ := strconv.ParseFloat(record[12], 64)
		y, _ := strconv.ParseFloat(record[13], 64)
		latitude, _ := strconv.ParseFloat(record[14], 64)
		longitude, _ := strconv.ParseFloat(record[15], 64)

		facility := models.Facility{
			LocationID:              record[0],
			Applicant:               record[1],
			FacilityType:            record[2],
			CNN:                     record[3],
			LocationDescription:     record[4],
			Address:                 record[5],
			BlockLot:                record[6],
			Block:                   record[7],
			Lot:                     record[8],
			Permit:                  record[9],
			Status:                  record[10],
			FoodItems:               record[11],
			X:                       x,
			Y:                       y,
			Latitude:                latitude,
			Longitude:               longitude,
			Schedule:                record[16],
			DaysHours:               record[17],
			NOISent:                 record[18],
			Approved:                record[19],
			Received:                record[20],
			PriorPermit:             record[21],
			ExpirationDate:          record[22],
			Location:                record[23],
			FirePreventionDistricts: record[24],
			PoliceDistricts:         record[25],
			SupervisorDistricts:     record[26],
			ZipCodes:                record[27],
			NeighborhoodsOld:        record[28],
		}
		if facility.Longitude > 0 || facility.Latitude > 0 {
			facilities = append(facilities, facility)
		}
	}

	return facilities, nil
}

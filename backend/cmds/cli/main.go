package main

import (
	"bufio"
	"context"
	"fmt"
	"food-trucks/packages/models"
	"food-trucks/packages/services"
	"food-trucks/packages/util/rdb"
	"food-trucks/packages/util/yaml"
	"os"
)

type CliConfig struct {
	Redis rdb.Config `yaml:"redis"`
}

func main() {
	svc := mustInit()
	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	for {
		fmt.Print("Enter Food Item to search facility: ")
		item, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		facilities, err := svc.GetByItem(ctx, item)
		if err != nil {
			panic(err)
		}
		for _, facility := range facilities {
			fmt.Println(facility.Applicant + " " + facility.LocationDescription)
		}
	}
}

func mustInit() *services.FacilitySvc {
	config, err := yaml.ParseYaml[CliConfig]("./configs/cli.yaml")
	if err != nil {
		panic(err)
	}
	facilityStore := rdb.NewEntityStore[string, models.Facility]("facility", 0, config.Redis).
		WithGetKey(models.GetFacilityKey)
	itemFacilityStore := rdb.NewSliceStore[string, models.Facility]("item", 0, config.Redis, facilityStore).
		WithGetKey(models.GetFacilityKey).WithGetScore(models.GetFacilityScore)
	geoFacilityStore := rdb.NewGeoStore[string, models.Facility]("geo", 0, config.Redis, facilityStore).
		WithGetKey(models.GetFacilityKey).WithGetLocation(models.GetFacilityLocation)
	facilitySvc := &services.FacilitySvc{
		FacilityStore:     facilityStore,
		ItemFacilityStore: itemFacilityStore,
		GeoFacilityStore:  geoFacilityStore,
	}
	err = facilitySvc.Seed("./configs/data.csv")
	if err != nil {
		panic(err)
	}
	return facilitySvc
}

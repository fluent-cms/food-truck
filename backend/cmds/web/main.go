package main

import (
	"fmt"
	"food-trucks/packages/controllers"
	"food-trucks/packages/models"
	"food-trucks/packages/services"
	"food-trucks/packages/util/irisbase"
	"food-trucks/packages/util/rdb"
	"food-trucks/packages/util/yaml"
)

type WebConfig struct {
	irisbase.AppConfig `yaml:"appConfig"`
	Redis              rdb.Config `yaml:"redis"`
}

type App struct {
	*irisbase.App
	WebConfig
}

type AppBuilder struct {
	WebConfig
}

func (b AppBuilder) BadRequest() []error {
	return []error{}
}

func (b AppBuilder) Services() []any {
	redisConfig := b.WebConfig.Redis
	fmt.Println("redisConfig:", redisConfig)
	facilityStore := rdb.NewEntityStore[string, models.Facility]("facility", 0, redisConfig).
		WithGetKey(models.GetFacilityKey)
	itemFacilityStore := rdb.NewSliceStore[string, models.Facility]("item", 0, redisConfig, facilityStore).
		WithGetKey(models.GetFacilityKey).WithGetScore(models.GetFacilityScore)
	geoFacilityStore := rdb.NewGeoStore[string, models.Facility]("geo", 0, redisConfig, facilityStore).
		WithGetKey(models.GetFacilityKey).WithGetLocation(models.GetFacilityLocation)
	facilitySvc := &services.FacilitySvc{
		FacilityStore:     facilityStore,
		ItemFacilityStore: itemFacilityStore,
		GeoFacilityStore:  geoFacilityStore,
	}
	err := facilitySvc.Seed("./configs/data.csv")
	if err != nil {
		panic(err)
	}
	return []any{facilitySvc}
}

func (b AppBuilder) Controller() map[string]any {
	return map[string]any{
		"/facilities/": new(controllers.FacilityCtl),
	}
}

func main() {
	config, err := yaml.ParseYaml[WebConfig]("./configs/web.yaml")
	if err != nil {
		panic(err)
	}
	app := &App{
		WebConfig: *config,
	}
	app.App = irisbase.NewIrisApp(config.AppConfig, AppBuilder{WebConfig: *config})
	app.Start()
}

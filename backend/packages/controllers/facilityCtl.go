package controllers

import (
	"food-trucks/packages/services"
	"github.com/kataras/iris/v12"
)

type FacilityCtl struct {
	C           iris.Context
	FacilitySvc *services.FacilitySvc
}

func (f FacilityCtl) GetCenter() any {
	return f.FacilitySvc.Center
}

func (f FacilityCtl) Get(qry struct {
	Lat    float64 `url:"lat"`
	Lon    float64 `url:"lon"`
	Radius float64 `url:"radius"`
}) any {
	items, err := f.FacilitySvc.GetByLocation(f.C.Request().Context(), qry.Lat, qry.Lon, qry.Radius)
	if err != nil {
		return err
	}
	return items
}

func (f FacilityCtl) GetBy(id string) any {
	item, err := f.FacilitySvc.GetByItem(f.C.Request().Context(), id)
	if err != nil {
		return err
	}
	return item
}

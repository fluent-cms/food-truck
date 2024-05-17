package rdb

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type TestEntity string
type EntityCache interface {
	Del(ctx context.Context, ids ...int) error
	Set(ctx context.Context, vals []TestEntity) error
	Get(ctx context.Context, keys []int) ([]TestEntity, error)
	GetFetch(ctx context.Context, keys []int, fetch func([]int) ([]TestEntity, error)) ([]TestEntity, error, error)
	GetFetchSet(ctx context.Context, keys []int, fetch func([]int) ([]TestEntity, error)) ([]TestEntity, error, error)
}

type SliceCache interface {
	SetSlice(ctx context.Context, sliceID any, items []TestEntity, left float64) error
	DelMember(ctx context.Context, sliceID any, memberKeys []int) error
	RevGetSlice(ctx context.Context, sliceID any, score *float64, count int) ([]TestEntity, float64, error)
	RevTruncate(ctx context.Context, sliceID any, score float64) error
}

type TimeSliceCache interface {
	GetLatest(ctx context.Context, sliceID any, count int, fetch func() ([]TestEntity, error)) ([]TestEntity, error, error)
	GetByTime(ctx context.Context, sliceID any, ts time.Time, count int, fetch func(ts time.Time, count int) ([]TestEntity, error)) ([]TestEntity, error)
	TruncateExpired(ctx context.Context, sliceID any) error
}

type TestService struct {
	entityCache    EntityCache
	sliceCache     SliceCache
	timeSliceCache TimeSliceCache
}

func NewTestService(entity EntityCache, slice SliceCache, timeSlice TimeSliceCache) *TestService {
	return &TestService{
		entityCache:    entity,
		sliceCache:     slice,
		timeSliceCache: timeSlice,
	}
}

// test entity, slice, timeSlice can  be injected to service
func TestNewTestService(t *testing.T) {
	entityStore := NewEntityStore[int, TestEntity]("testEntities", time.Minute*60, getConfig())
	sliceStore := NewSliceStore[int, TestEntity]("testEntities", time.Minute*60, getConfig(), entityStore)
	timeSliceStore := NewTimeSliceStore[int, TestEntity]("testEntities", time.Minute*60, time.Minute, getConfig(), entityStore)
	service := NewTestService(entityStore, sliceStore, timeSliceStore)
	fmt.Println(service)
}

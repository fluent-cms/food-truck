package rdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/constraints"
	"strings"
	"time"
)

type GeoStore[K constraints.Ordered, Entity any] struct {
	Config
	namespace    string
	duration     time.Duration
	client       redis.UniversalClient
	entityStore  Cacheable[K, Entity]
	getMemberKey func(Entity) K
	getLocation  func(Entity) (float64, float64)
}

func NewGeoStore[MemberKey constraints.Ordered, Entity any](
	namespace string,
	duration time.Duration,
	config Config,
	entityStore Cacheable[MemberKey, Entity],
) *GeoStore[MemberKey, Entity] {
	return &GeoStore[MemberKey, Entity]{
		Config:      config,
		client:      redis.NewUniversalClient(&redis.UniversalOptions{Addrs: strings.Split(config.Addr, ",")}),
		namespace:   namespace,
		entityStore: entityStore,
		duration:    duration,
	}
}

func (s *GeoStore[K, V]) WithGetKey(f func(V) K) *GeoStore[K, V] {
	s.getMemberKey = f
	return s
}

func (s *GeoStore[K, V]) WithGetLocation(f func(V) (lat float64, lon float64)) *GeoStore[K, V] {
	s.getLocation = f
	return s
}

func (s *GeoStore[K, V]) Add(ctx context.Context, item V) error {
	lat, lon := s.getLocation(item)
	key := fmt.Sprintf("%v", s.getMemberKey(item))
	_, err := s.client.GeoAdd(ctx, s.namespace, &redis.GeoLocation{
		Name:      key,
		Longitude: lon,
		Latitude:  lat,
	}).Result()
	return err
}

func (s *GeoStore[K, V]) Get(ctx context.Context, lat float64, lon float64, radius float64) ([]V, error) {
	res, err := s.client.GeoRadius(ctx, s.namespace, lon, lat, &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      "km",
		WithCoord: true,
		WithDist:  true,
	}).Result()
	if err != nil {
		return nil, err
	}
	var keys []K
	for _, re := range res {
		k, err := keyFromStar[K](re.Name)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return s.entityStore.Get(ctx, keys)
}

package rdb

import (
	"context"
	"fmt"
	"food-trucks/packages/util/annotate"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
	"strings"
	"time"
)

type SliceStore[K constraints.Ordered, Entity any] struct {
	Config
	namespace    string
	duration     time.Duration
	client       redis.UniversalClient
	entityStore  Cacheable[K, Entity]
	getMemberKey func(Entity) K
	getScore     func(Entity) float64
}

func NewSliceStore[MemberKey constraints.Ordered, Entity any](
	namespace string,
	duration time.Duration,
	config Config,
	entityStore Cacheable[MemberKey, Entity],
) *SliceStore[MemberKey, Entity] {
	return &SliceStore[MemberKey, Entity]{
		Config:      config,
		client:      redis.NewUniversalClient(&redis.UniversalOptions{Addrs: strings.Split(config.Addr, ",")}),
		namespace:   namespace,
		entityStore: entityStore,
		duration:    duration,
	}
}

func (s *SliceStore[K, V]) WithGetKey(f func(V) K) *SliceStore[K, V] {
	s.getMemberKey = f
	return s
}

func (s *SliceStore[K, V]) WithGetScore(f func(V) float64) *SliceStore[K, V] {
	s.getScore = f
	return s
}

func (s *SliceStore[K, V]) DelSlice(ctx context.Context, sliceID any) error {
	key := s.sliceKey(sliceID)
	_, err := s.client.Del(ctx, key).Result()
	return err
}

func (s *SliceStore[K, V]) DelMember(ctx context.Context, sliceID any, memberKeys []K) error {
	key := s.sliceKey(sliceID)
	_, err := s.client.ZRem(ctx, key, lo.ToAnySlice(memberKeys)...).Result()
	return err
}

func (s *SliceStore[K, Entity]) AddMem(ctx context.Context, sliceID any, items []Entity) error {
	key := s.sliceKey(sliceID)
	if err := s.entityStore.Set(ctx, items); err != nil {
		return annotate.Error(err)
	}
	members := make([]redis.Z, len(items))
	for i, item := range items {
		members[i] = redis.Z{
			Score:  s.getScore(item),
			Member: s.getMemberKey(item),
		}
	}
	_, err := s.client.ZAdd(ctx, key, members...).Result()
	return err
}

func (s *SliceStore[K, Entity]) GetAllMemberEntities(ctx context.Context, sliceID any) ([]Entity, error) {
	option := &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}
	key := s.sliceKey(sliceID)
	p := s.client.Pipeline()
	memCmd := p.ZRevRangeByScore(ctx, key, option)
	if _, err := p.Exec(ctx); IgnoreNoKey(err) != nil {
		return nil, annotate.Error(err)
	}
	items, err := s.getEntities(ctx, memCmd.Val())
	if err != nil {
		return nil, annotate.Error(err)
	}

	return items, nil
}

func (s *SliceStore[K, V]) getEntities(ctx context.Context, members []string) ([]V, error) {
	var memberKeys []K
	for _, mem := range members {
		memberK, err := keyFromStar[K](mem)
		if err != nil {
			return nil, annotate.Error(err)
		}
		memberKeys = append(memberKeys, memberK)
	}
	return s.entityStore.Get(ctx, memberKeys)
}

func (s *SliceStore[K, V]) sliceKey(sliceID any) string {
	return fmt.Sprintf("%s:%s:%v", s.Prefix, s.namespace, sliceID)
}

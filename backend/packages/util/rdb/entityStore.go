package rdb

import (
	"context"
	"errors"
	"fmt"
	"food-trucks/packages/util/singleflight"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
	"time"
)

type EntityStore[K constraints.Ordered, V any] struct {
	client         *Client[K]
	namespace      string
	entityDuration time.Duration
	getKey         func(V) K
	single         *singleflight.Group[string, []V]
}

type gGetResult[K constraints.Ordered, V any] struct {
	Values []V
	Missed []K
}

func NewEntityStore[K constraints.Ordered, V any](namespace string, duration time.Duration, config Config) *EntityStore[K, V] {
	return &EntityStore[K, V]{
		single:         &singleflight.Group[string, []V]{},
		client:         NewClient[K](config),
		namespace:      namespace,
		entityDuration: duration,
	}
}

func (c *EntityStore[K, V]) WithGetKey(f func(V) K) *EntityStore[K, V] {
	c.getKey = f
	return c
}

func (c *EntityStore[K, V]) Set(ctx context.Context, vals []V) error {
	items, err := c.toStrEntity(vals)
	if err != nil {
		return err
	}
	return c.client.mSet(ctx, c.namespace, c.entityDuration, items)
}

func (c *EntityStore[K, V]) Del(ctx context.Context, ids ...K) error {
	return c.client.mDel(ctx, c.namespace, ids)
}

func (c *EntityStore[K, V]) Get(ctx context.Context, keys []K) ([]V, error) {
	return c.getFetchSet(ctx, keys, false, nil)
}

func (c *EntityStore[K, V]) GetFetch(ctx context.Context, keys []K,
	fetch func([]K) ([]V, error),
) ([]V, error, error) {
	singleKey := fmt.Sprintf("%v", keys)
	vals, err, _ := c.single.Do(singleKey, func() ([]V, error) {
		vals, err := c.getFetchSet(ctx, keys, false, fetch)
		return vals, err
	})
	return errOrWarning(vals, err)
}

func (c *EntityStore[K, V]) GetFetchSet(ctx context.Context, keys []K,
	fetch func([]K) ([]V, error),
) ([]V, error, error) {
	singleKey := fmt.Sprintf("S:%v", keys)
	vals, err, _ := c.single.Do(singleKey, func() ([]V, error) {
		return c.getFetchSet(ctx, keys, true, fetch)
	})
	return errOrWarning(vals, err)
}

func (c *EntityStore[K, V]) toStrEntity(values []V) ([]lo.Entry[K, string], error) {
	if c.getKey == nil {
		return nil, errors.New("getKey not set")
	}
	var ret []lo.Entry[K, string]
	for _, value := range values {
		str, err := toStr(value)
		if err != nil {
			return nil, err
		}
		ret = append(ret, lo.Entry[K, string]{
			Key:   c.getKey(value),
			Value: str,
		})
	}
	return ret, nil
}

func (c *EntityStore[K, V]) get(ctx context.Context, keys []K) (gGetResult[K, V], error) {
	ret := gGetResult[K, V]{}
	items, missed, err := c.client.mGet(ctx, c.namespace, keys)
	if err != nil {
		return ret, nil
	}
	vals, err := mFromStr[V](items)
	if err != nil {
		return ret, nil
	}
	ret.Missed = missed
	ret.Values = vals
	return ret, nil
}

func (c *EntityStore[K, V]) getFetchSet(ctx context.Context, keys []K, cacheFetchResult bool,
	fetch func([]K) ([]V, error),
) ([]V, error) {
	if c.getKey == nil {
		return nil, errors.New("getKey not set")
	}

	ret, err := c.get(ctx, keys)
	if err != nil {
		return nil, err
	}

	var warning error
	if len(ret.Missed) > 0 && fetch != nil {
		items, err := fetch(ret.Missed)
		if err != nil {
			return nil, err
		}
		ret.Values = append(ret.Values, items...)

		if cacheFetchResult {
			warning = wrapSetCacheError(c.Set(ctx, items))
		}
	}
	return c.sort(keys, ret.Values), warning
}

func (c *EntityStore[K, V]) sort(keys []K, vals []V) []V {
	return lo.Map(keys, func(k K, index int) V {
		item, _ := lo.Find(vals, func(v V) bool {
			return k == c.getKey(v)
		})
		return item
	})
}

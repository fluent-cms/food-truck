package rdb

import (
	"context"
	"golang.org/x/exp/constraints"
)

type Cacheable[K constraints.Ordered, V any] interface {
	Set(ctx context.Context, vals []V) error
	Get(ctx context.Context, keys []K) ([]V, error)
}

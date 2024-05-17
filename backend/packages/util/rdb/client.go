package rdb

import (
	"context"
	"fmt"
	"food-trucks/packages/util/safeslice"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

type Config struct {
	Addr            string `yaml:"addr"`
	Prefix          string `yaml:"prefix"`
	HashtagPosition int    `yaml:"hashtagPosition"`
	Gzip            bool   `yaml:"gzip"`
	Enabled         bool   `yaml:"enabled"`
}

type Client[K constraints.Ordered] struct {
	Config
	Client redis.UniversalClient
}

func NewClient[K constraints.Ordered](config Config) *Client[K] {
	return &Client[K]{
		Config: config,
		Client: redis.NewUniversalClient(&redis.UniversalOptions{Addrs: strings.Split(config.Addr, ",")}),
	}
}

func (c *Client[K]) makeKey(namespace string, id K) string {
	return fmt.Sprintf("%s:%s:%v", c.Prefix, namespace, id)
}

func (c *Client[K]) mDel(ctx context.Context, namespace string, ids []K) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, entities := range c.tagKeys(namespace, ids) {
		g.Go(func() error {
			keys := lo.Map(entities, func(item lo.Entry[string, K], index int) string {
				return item.Key
			})
			_, err := c.Client.Del(ctx, keys...).Result()
			return err
		})
	}
	return g.Wait()
}

func (c *Client[K]) del(ctx context.Context, namespace string, id K) error {
	_, err := c.Client.Del(ctx, c.makeKey(namespace, id)).Result()
	return err
}

func (c *Client[K]) mGet(ctx context.Context, namespace string, ids []K) ([]string, []K, error) {
	g, ctx := errgroup.WithContext(ctx)
	ret := safeslice.NewSafeSlice[string]()
	missed := safeslice.NewSafeSlice[K]()
	for _, entities := range c.tagKeys(namespace, ids) {
		g.Go(func() error {
			keys := lo.Map(entities, func(item lo.Entry[string, K], index int) string {
				return item.Key
			})

			res, err := c.Client.MGet(ctx, keys...).Result()
			if err != nil {
				return err
			}

			var vals []string
			for i, re := range res {
				if re == nil {
					missed.Append(entities[i].Value)
				} else {
					vals = append(vals, re.(string))
				}
			}

			if c.Gzip {
				var err error
				vals, err = munzip(vals)
				if err != nil {
					return err
				}
			}
			ret.Append(vals...)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}
	return ret.All(), missed.All(), nil
}

func (c *Client[K]) get(ctx context.Context, namespace string, id K) (string, error) {
	str, err := c.Client.Get(ctx, c.makeKey(namespace, id)).Result()
	if err != nil {
		return "", err
	}
	if c.Gzip {
		str, err = unzip(str)
		if err != nil {
			return "", err
		}
	}
	return str, nil
}

func (c *Client[K]) mSet(ctx context.Context, namespace string, duration time.Duration, entities []lo.Entry[K, string]) error {
	g, ctx := errgroup.WithContext(ctx)
	if c.Gzip {
		for i, entity := range entities {
			var err error
			entity.Value, err = zip(entity.Value)
			if err != nil {
				return err
			}
			entities[i] = entity
		}
	}
	for _, items := range c.tagEntities(namespace, entities) {
		g.Go(func() error {
			p := c.Client.Pipeline()
			for _, entity := range items {
				p.Set(ctx, entity.Key, entity.Value, duration)
			}
			_, err := p.Exec(ctx)
			return err
		})
	}
	return g.Wait()
}

func (c *Client[K]) set(ctx context.Context, namespace string, id K, value string, expiration time.Duration) error {
	var err error
	if c.Gzip {
		value, err = zip(value)
		if err != nil {
			return err
		}
	}
	_, err = c.Client.Set(ctx, c.makeKey(namespace, id), value, expiration).Result()
	return err
}

type TagKey struct {
	Tag string
	Key string
}

func (c *Client[K]) tag(namespace string, id K) TagKey {
	str := c.makeKey(namespace, id)
	ret := TagKey{Key: str}
	if len(str) > c.HashtagPosition {
		pos := len(str) - c.HashtagPosition
		ret.Tag = "{" + str[:pos] + "}"
		ret.Key = ret.Tag + str[pos:]
	}
	return ret
}

func (c *Client[K]) tagEntities(namespace string, entities []lo.Entry[K, string]) map[string][]lo.Entry[string, string] {
	ret := make(map[string][]lo.Entry[string, string])
	for _, entity := range entities {
		tag := c.tag(namespace, entity.Key)
		ret[tag.Tag] = append(ret[tag.Tag], lo.Entry[string, string]{Key: tag.Key, Value: entity.Value})
	}
	return ret
}

func (c *Client[K]) tagKeys(namespace string, ids []K) map[string][]lo.Entry[string, K] {
	ret := make(map[string][]lo.Entry[string, K])
	for _, id := range ids {
		tag := c.tag(namespace, id)
		ret[tag.Tag] = append(ret[tag.Tag], lo.Entry[string, K]{
			Key:   tag.Key,
			Value: id,
		})
	}
	return ret
}

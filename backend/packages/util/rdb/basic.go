package rdb

import (
	"context"
	"time"
)

/*
performance is not optimized and not scalable, just for simple use cases;
for advance usage, e.g. multiple redis shard and batch operation for better performance, use entityStore
*/
var client *Client[string]

func Init(config Config) {
	client = NewClient[string](config)
}

func GetStr(ctx context.Context, namespace string, key string) (string, error) {
	return client.get(ctx, namespace, key)
}

func SetStr(ctx context.Context, namespace string, key string, value string, expiration time.Duration) error {
	return client.set(ctx, namespace, key, value, expiration)
}

func Del(ctx context.Context, namespace string, key string) error {
	return client.del(ctx, namespace, key)
}

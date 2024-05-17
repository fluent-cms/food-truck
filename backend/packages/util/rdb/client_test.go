package rdb

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"testing"
	"time"
)

// start redis from local
// docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
// docker exec -it redis-stack redis-cli

func getConfig() Config {
	return Config{
		Addr:            "127.0.0.1:6379",
		Prefix:          "Test",
		HashtagPosition: 3,
		Gzip:            false,
		Enabled:         true,
	}
}
func getIntKeyClient() *Client[int] {
	return NewClient[int](getConfig())
}

func TestMSet(t *testing.T) {
	c := getIntKeyClient()
	err := c.mSet(context.Background(), "posts", time.Hour, []lo.Entry[int, string]{
		{
			Key:   10001,
			Value: "v10001",
		},
		{
			Key:   10002,
			Value: "v10002",
		},
		{
			Key:   10003,
			Value: "v10003",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestMGet(t *testing.T) {
	c := getIntKeyClient()
	vals, missed, err := c.mGet(context.Background(), "posts", []int{10001, 10008})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(vals, missed)
}

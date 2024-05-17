package rdb

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"testing"
	"time"
)

type EntityStorePost struct {
	ID    int
	Title string
}

func EntityStorePostID(p EntityStorePost) int {
	return p.ID
}

var TestEntityStore = "TestEntityStore"

func TestEntityStore_Set(t *testing.T) {
	ctx := context.Background()
	entityStore := NewEntityStore[int, EntityStorePost](TestEntityStore, time.Hour, getConfig()).
		WithGetKey(EntityStorePostID)
	if err := entityStore.Set(ctx, fakeEntityStorePost(10)); err != nil {
		t.Fatal(err)
	}
}

func TestEntityStore_Del(t *testing.T) {
	ctx := context.Background()
	entityStore := NewEntityStore[int, EntityStorePost](TestEntityStore, time.Hour, getConfig()).
		WithGetKey(EntityStorePostID)
	if err := entityStore.Del(ctx, 1001, 1002, 1003, 1004, 1005, 1007); err != nil {
		t.Fatal(err)
	}
}

func TestEntityStore_Get(t *testing.T) {
	ctx := context.Background()
	entityStore := NewEntityStore[int, EntityStorePost](TestEntityStore, time.Hour, getConfig()).
		WithGetKey(EntityStorePostID)
	if items, err := entityStore.Get(ctx, fakeEntityStoreIDs(10)); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(items)
	}
}

func TestEntityStore_GetFetch(t *testing.T) {
	ctx := context.Background()
	entityStore := NewEntityStore[int, EntityStorePost](TestEntityStore, time.Hour, getConfig()).
		WithGetKey(EntityStorePostID)
	if items, err, _ := entityStore.GetFetch(ctx, fakeEntityStoreIDs(10), fakeEntityPostByIDs); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(items)
	}
}

func TestEntityStore_GetFetchSet(t *testing.T) {
	ctx := context.Background()
	entityStore := NewEntityStore[int, EntityStorePost](TestEntityStore, time.Hour, getConfig()).
		WithGetKey(EntityStorePostID)
	if items, err, _ := entityStore.GetFetchSet(ctx, fakeEntityStoreIDs(10), fakeEntityPostByIDs); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(items)
	}
}

func fakeEntityPostByIDs(ints []int) ([]EntityStorePost, error) {
	return lo.Map(ints, func(item int, index int) EntityStorePost {
		return EntityStorePost{
			ID:    item,
			Title: fmt.Sprintf("fetch again %v", item),
		}
	}), nil

}

func fakeEntityStorePost(count int) []EntityStorePost {
	var ret []EntityStorePost
	for i := 0; i < count; i++ {
		ret = append(ret, EntityStorePost{
			ID:    1000 + i,
			Title: fmt.Sprintf("Post %v", i),
		},
		)
	}
	return ret
}

func fakeEntityStoreIDs(count int) []int {
	var ret []int
	for i := 0; i < count; i++ {
		ret = append(ret, 1000+i)
	}
	return ret
}

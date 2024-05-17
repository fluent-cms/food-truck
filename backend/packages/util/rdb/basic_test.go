package rdb

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestSetStr(t *testing.T) {
	Init(getConfig())
	err := SetStr(context.Background(), "user", "jk", "2002 JK ", time.Hour)
	if err != nil {
		t.Fatal()
	}
}

func TestGetStr(t *testing.T) {
	Init(getConfig())
	s, err := GetStr(context.Background(), "user", "jk")
	if IgnoreNoKey(err) != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
}

func TestDel(t *testing.T) {
	Init(getConfig())
	err := Del(context.Background(), "user", "jk")
	if err != nil {
		t.Fatal()
	}
}

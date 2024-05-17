package rdb

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/constraints"
	"io/ioutil"
	"strconv"
	"strings"
)

func zip(input string) (string, error) {
	var compressed bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressed)
	_, err := gzipWriter.Write([]byte(input))
	if err != nil {
		return "", fmt.Errorf("error gzip write %w", err)
	}
	err = gzipWriter.Close()
	if err != nil {
		return "", fmt.Errorf("error gzip write %w", err)
	}
	return string(compressed.Bytes()), nil
}

func munzip(vals []string) ([]string, error) {
	var ret []string
	var err error
	for _, val := range vals {
		val, err = unzip(val)
		if err != nil {
			return nil, err
		}
		ret = append(ret, val)
	}
	return ret, nil
}
func unzip(compressedData string) (string, error) {
	compressedReader := bytes.NewReader([]byte(compressedData))
	gzipReader, err := gzip.NewReader(compressedReader)
	if err != nil {
		return "", fmt.Errorf("error creating gzip reader: %w", err)
	}
	decompressed, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return "", fmt.Errorf("error decompressing data: %w", err)
	}
	err = gzipReader.Close()
	if err != nil {
		return "", fmt.Errorf("error closing gzip reader: %w", err)
	}
	return string(decompressed), nil
}

func fromByte[V any](bs []byte) (V, error) {
	var v V
	var anyV any
	anyV = v
	switch anyV.(type) {
	case string:
		anyV = string(bs)
		v = anyV.(V)
		return v, nil
	case int:
		anyV = int(binary.LittleEndian.Uint64(bs))
		v = anyV.(V)
		return v, nil
	case uint:
		anyV = uint(binary.LittleEndian.Uint64(bs))
		v = anyV.(V)
		return v, nil
	case int64:
		anyV = int64(binary.LittleEndian.Uint64(bs))
		v = anyV.(V)
		return v, nil
	case uint64:
		anyV = binary.LittleEndian.Uint64(bs)
		v = anyV.(V)
		return v, nil
	default:
		err := json.Unmarshal(bs, &v)
		return v, err
	}
}

func toBytes(k any) ([]byte, error) {
	switch i := k.(type) {
	case string:
		return []byte(i), nil
	case int:
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, uint64(i))
		return bs, nil
	case uint:
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, uint64(i))
		return bs, nil
	case int64:
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, uint64(i))
		return bs, nil
	case uint64:
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, i)
		return bs, nil
	default:
		return json.Marshal(k)
	}
}
func fromStr[T any](str string) (T, error) {
	return fromByte[T]([]byte(str))
}

func mFromStr[T any](strs []string) ([]T, error) {
	var ret []T
	for _, str := range strs {
		s, err := fromStr[T](str)
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}

func toStr(v any) (string, error) {
	bs, err := toBytes(v)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func keyFromStar[K constraints.Ordered](i string) (K, error) {
	var err error
	var k K
	var a any = k
	switch a.(type) {
	case int:
		a, err = strconv.Atoi(i)
		if err != nil {
			return k, fmt.Errorf("unsupported type")
		}
		return a.(K), nil
	case int64:
		a, err = strconv.ParseInt(i, 10, 64)
		if err != nil {
			return k, fmt.Errorf("unsupported type")
		}
		return a.(K), nil
	case string:
		a = i
		return a.(K), nil
	default:
		return k, fmt.Errorf("unsupported type")
	}
}

func IgnoreNoKey(err error) error {
	switch {
	case err == nil, errors.Is(err, redis.Nil), strings.Contains(err.Error(), "key that doesn't exist"):
		return nil
	default:
		return err
	}
}

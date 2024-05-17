package rdb

import (
	"errors"
	"fmt"
)

// SetCacheError most of the time we can ignore this error, this error often return as warning.
var SetCacheError = errors.New("fail to set cache")

func wrapSetCacheError(err error) error {
	if err != nil {
		return fmt.Errorf("%w, %v", SetCacheError, err)
	}
	return nil
}

func errOrWarning[V any](vals V, err error) (V, error, error) {
	if errors.Is(err, SetCacheError) {
		return vals, nil, err
	} else {
		return vals, err, nil
	}
}

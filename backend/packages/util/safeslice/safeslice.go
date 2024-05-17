package safeslice

import "sync"

type SafeSlice[T any] struct {
	sync.Mutex
	items []T
}

func NewSafeSlice[T any]() *SafeSlice[T] {
	return &SafeSlice[T]{}
}

func (s *SafeSlice[T]) Append(items ...T) {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		s.items = append(s.items, item)
	}
}

func (s *SafeSlice[T]) Get(index int) T {
	s.Lock()
	defer s.Unlock()
	return s.items[index]
}

func (s *SafeSlice[T]) All() []T {
	s.Lock()
	defer s.Unlock()
	return s.items
}

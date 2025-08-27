package services

import "sync"

type Remeber[T comparable] interface {
	Remember(T) bool
	Clear()
}

type remember[T comparable] struct {
	values map[T]struct{}
	mu     sync.RWMutex
}

func NewRemember[T comparable]() Remeber[T] {
	return &remember[T]{
		values: make(map[T]struct{}),
	}
}

func (r *remember[T]) Remember(val T) bool {
	r.mu.RLock()

	if _, ok := r.values[val]; ok {
		r.mu.RUnlock()
		return true
	}

	r.mu.RUnlock()
	r.mu.Lock()
	defer r.mu.Unlock()

	r.values[val] = struct{}{}

	return false
}

func (r *remember[T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k := range r.values {
		delete(r.values, k)
	}
}

package store

import (
	lru "github.com/hashicorp/golang-lru"
)

type Storer interface {
	// Get retrieves a value from the store.
	Get(key string) (any, bool)
	// Add stores a value in the store.
	Add(key, value any) bool
	// Delete removes a value from the store.
	Delete(key string) bool
}

type store struct {
	lru *lru.Cache
}

// New creates a new store with the given size.
func New(size int) (Storer, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}

	return &store{
		lru: cache,
	}, nil
}

func (s *store) Get(key string) (any, bool) {
	return s.lru.Get(key)
}

func (s *store) Add(key, value any) bool {
	return s.lru.Add(key, value)
}

func (s *store) Delete(key string) bool {
	return s.lru.Remove(key)
}

// Package store implements an in-memory key-value store with optional key expiry.
package store

import (
	"cmp"
	"errors"
	"slices"
	"time"
)

const NoExpiry time.Duration = -1

var ErrNotFound = errors.New("value not found")

type Store struct {
	clock       func() time.Time
	values      map[string]string
	expirations map[string]time.Time
}

func New() *Store {
	return &Store{
		clock:       time.Now,
		values:      map[string]string{},
		expirations: map[string]time.Time{},
	}
}

func (store *Store) Set(key, value string) {
	store.values[key] = value
	delete(store.expirations, key)
}

func (store *Store) SetWithTTL(key, value string, ttl time.Duration) {
	store.Set(key, value)
	store.expirations[key] = store.clock().Add(ttl)
}

func (store *Store) Get(key string) (string, error) {
	store.deleteIfExpired(key)

	value, ok := store.values[key]

	if !ok {
		return "", ErrNotFound
	}

	return value, nil
}

func (store *Store) Delete(key string) bool {
	store.deleteIfExpired(key)
	_, ok := store.values[key]

	if !ok {
		return false
	}

	delete(store.values, key)
	delete(store.expirations, key)

	return true
}

func (store *Store) Keys() []string {
	for key := range store.values {
		store.deleteIfExpired(key)
	}

	keys := make([]string, 0, len(store.values))

	for key := range store.values {
		keys = append(keys, key)
	}

	slices.SortFunc(keys, func(a, b string) int {
		return cmp.Compare(a, b)
	})

	return keys
}

func (store *Store) TTL(key string) (time.Duration, error) {
	store.deleteIfExpired(key)

	if _, exists := store.values[key]; !exists {
		return 0, ErrNotFound
	}

	expiresAt, hasExpiry := store.expirations[key]

	if !hasExpiry {
		return NoExpiry, nil
	}

	return expiresAt.Sub(store.clock()), nil
}

func (store *Store) deleteIfExpired(key string) {
	ttl, hasLifetime := store.expirations[key]

	if !hasLifetime {
		return
	}

	if store.clock().After(ttl) {
		delete(store.values, key)
		delete(store.expirations, key)
	}
}

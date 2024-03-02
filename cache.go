package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	mux     sync.Mutex
	expired ExpirationStrategy[V]
	values  map[K]*item[V]
}

type ExpirationStrategy[V any] func(i *item[V]) bool

type item[V any] struct {
	tm    time.Time
	value V
}

func New[K comparable, V any](expirationStrategy ExpirationStrategy[V]) *Cache[K, V] {
	return &Cache[K, V]{
		expired: expirationStrategy,
		values:  make(map[K]*item[V]),
	}
}

func (c *Cache[K, V]) Peek(k K) (V, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	i, ok := c.values[k]
	if !ok {
		return *new(V), false
	}
	return i.value, ok
}

func (c *Cache[K, V]) Get(k K, get func(k K) (V, error)) (V, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	i, ok := c.values[k]
	if ok && c.expired(i) {
		delete(c.values, k)
		ok = false
	}
	if !ok {
		n, err := get(k)
		if err != nil {
			return n, err
		}
		i = &item[V]{
			tm:    time.Now(),
			value: n,
		}
		c.values[k] = i
	}
	return i.value, nil
}

func (c *Cache[K, V]) Put(k K, v V) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.values[k] = &item[V]{
		tm:    time.Now(),
		value: v,
	}
}

func (c *Cache[K, V]) Remove(keys ...K) {
	c.mux.Lock()
	defer c.mux.Unlock()
	for _, k := range keys {
		delete(c.values, k)
	}
}

// NoExpiration that items never expired
func NoExpiration[V any]() ExpirationStrategy[V] {
	return func(i *item[V]) bool {
		return false
	}
}

// ExpireAfter item expired after N seconds after a put to the cache
func ExpireAfter[V any](dt time.Duration) ExpirationStrategy[V] {
	return func(i *item[V]) bool {
		return time.Since(i.tm) > dt
	}
}

// LastAccess item expired after N second since last access
func LastAccess[V any](dt time.Duration) ExpirationStrategy[V] {
	return func(i *item[V]) bool {
		now := time.Now()
		if now.Sub(i.tm) > dt {
			return true
		}
		i.tm = now
		return false
	}
}

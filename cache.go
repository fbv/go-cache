package cache

import (
	"time"
)

type Cache[K comparable, V any] struct {
	expired ExpirationStrategy[V]
	get     func(K) (V, error)
	values  map[K]*item[V]
}

type ExpirationStrategy[V any] func(i *item[V]) bool

type item[V any] struct {
	tm    time.Time
	value V
}

func New[K comparable, V any](expirationStrategy ExpirationStrategy[V], get func(K) (V, error)) *Cache[K, V] {
	return &Cache[K, V]{
		expired: expirationStrategy,
		get:     get,
		values:  make(map[K]*item[V]),
	}
}

func (c *Cache[K, V]) Get(k K) (V, error) {
	i, ok := c.values[k]
	if !ok || c.expired(i) {
		n, err := c.get(k)
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

func NoExpiration[V any]() ExpirationStrategy[V] {
	return func(i *item[V]) bool {
		return false
	}
}

func ExpireAfter[V any](dt time.Duration) ExpirationStrategy[V] {
	return func(i *item[V]) bool {
		return time.Since(i.tm) > dt
	}
}

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

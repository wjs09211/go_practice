// Package localcache
package localcache

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotExist = errors.New("key not exist")
)

// Cache interface implement by localCache
type Cache interface {
	// Get value from key
	Get(key string) (any, error)
	// Set value by key
	Set(key string, value any, expiredTime time.Duration) error
}

// New localCache
func New() Cache {
	return &localCache{make(map[string]*cacheData), sync.RWMutex{}}
}

type localCache struct {
	hash  map[string]*cacheData
	mutex sync.RWMutex
}

type cacheData struct {
	data      any
	expiredAt time.Time
}

func (lc *localCache) Get(key string) (any, error) {
	lc.mutex.RLock()
	defer lc.mutex.RUnlock()
	val, ok := lc.hash[key]

	if !ok {
		return nil, ErrKeyNotExist
	}

	if val.expiredAt.Before(time.Now()) {
		delete(lc.hash, key)
		return nil, ErrKeyNotExist
	}

	return val.data, nil

}

func (lc *localCache) Set(key string, value any, expiredTime time.Duration) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	lc.hash[key] = &cacheData{value, time.Now().Add(expiredTime)}
	return nil
}

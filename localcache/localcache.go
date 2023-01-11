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
	return &localCache{make(map[string]*cacheData), sync.Mutex{}}
}

type localCache struct {
	hash  map[string]*cacheData
	mutex sync.Mutex
}

type cacheData struct {
	data      any
	expiredAt time.Time
}

func (obj *localCache) Get(key string) (any, error) {
	if val, ok := obj.hash[key]; ok {
		if val.expiredAt.Before(time.Now()) {
			return nil, nil
		}
		return val.data, nil
	}
	return nil, ErrKeyNotExist
}

func (obj *localCache) Set(key string, value any, expiredTime time.Duration) error {
	obj.mutex.Lock()
	obj.hash[key] = &cacheData{value, time.Now().Add(expiredTime)}
	defer obj.mutex.Unlock()
	return nil
}

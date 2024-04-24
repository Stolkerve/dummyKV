package cache

import (
	"sync"
	"time"
)

const (
	NO_EXPIRATION time.Duration = -1
)

type Object struct {
	data       interface{}
	expiration int64
}

type Cache struct {
	objects map[string]Object
	rwm     sync.RWMutex
}

func NewCache() Cache {
	return Cache{
		objects: make(map[string]Object),
	}
}

func (c *Cache) Get(name string) (interface{}, bool) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()

	obj, ok := c.objects[name]
	if !ok {
		return nil, false
	}

	if obj.expiration > 0 {
		if time.Now().UnixMilli() > obj.expiration {
			delete(c.objects, name)
			return nil, false
		}
	}

	return obj.data, ok
}

// Set or update a value
func (c *Cache) Set(name string, value interface{}, expiration time.Duration) interface{} {
	c.rwm.Lock()
	defer c.rwm.Unlock()

	e := int64(NO_EXPIRATION)
	if expiration > 0 {
		e = int64(time.Now().Add(expiration).UnixMilli())
	}

	oldValue := c.objects[name]
	c.objects[name] = Object{
		data:       value,
		expiration: e,
	}

	return oldValue
}

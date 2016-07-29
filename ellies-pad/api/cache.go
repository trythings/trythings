package api

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Cache struct {
	objs map[string]interface{}
}

func (c *Cache) Get(key *datastore.Key) interface{} {
	if c == nil {
		return nil
	}

	return c.objs[key.String()]
}

func (c *Cache) Set(key *datastore.Key, value interface{}) {
	if c == nil {
		return
	}

	c.objs[key.String()] = value
}

func NewCacheContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, cacheKey, &Cache{
		objs: map[string]interface{}{},
	})
}

func CacheFromContext(ctx context.Context) *Cache {
	c, ok := ctx.Value(cacheKey).(*Cache)
	if !ok {
		return nil
	}

	return c
}

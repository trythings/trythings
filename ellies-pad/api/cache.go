package api

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Cache struct {
	objs           map[string]interface{}
	spaceIsVisible map[string]bool
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

// TODO: Looks like each service should maintain its own per-request cache.
// TODO: When we support adding Users to a Space, we'll need to make sure we invalidate these cache entries.
func (c *Cache) IsVisible(sp *Space) (bool, bool) {
	if c == nil {
		return false, false
	}

	isVisible, ok := c.spaceIsVisible[sp.ID]
	return isVisible, ok
}

func (c *Cache) SetIsVisible(sp *Space, isVisible bool) {
	if c == nil {
		return
	}

	c.spaceIsVisible[sp.ID] = isVisible
}

func NewCacheContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, cacheKey, &Cache{
		objs:           map[string]interface{}{},
		spaceIsVisible: map[string]bool{},
	})
}

func CacheFromContext(ctx context.Context) *Cache {
	c, ok := ctx.Value(cacheKey).(*Cache)
	if !ok {
		return nil
	}

	return c
}

package api

import "golang.org/x/net/context"

type Cache struct {
	objs map[string]interface{}
}

func (c *Cache) Get(id string) interface{} {
	if c == nil {
		return nil
	}

	return c.objs[id]
}

func (c *Cache) Set(id string, value interface{}) {
	if c == nil {
		return
	}

	c.objs[id] = value
}

func NewPerRequestCacheContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, perRequestCacheKey, &Cache{
		objs: map[string]interface{}{},
	})
}

func CacheFromContext(ctx context.Context) *Cache {
	c, ok := ctx.Value(perRequestCacheKey).(*Cache)
	if !ok {
		return nil
	}

	return c
}

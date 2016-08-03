package api

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Cache struct {
	objs           map[string]interface{}
	spaceIsVisible map[string]bool
	// SpaceID -> query -> results
	searchResults map[string]map[string][]*Task
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

// TODO: Figure out which updates invalidate a search's results (assuming mutations may happen after fetches).
func (c *Cache) SearchResults(sp *Space, query string) ([]*Task, bool) {
	if c == nil {
		return nil, false
	}

	perQuery := c.searchResults[sp.ID]
	if perQuery == nil {
		return nil, false
	}

	ts, ok := perQuery[query]
	return ts, ok
}

func (c *Cache) SetSearchResults(sp *Space, query string, results []*Task) {
	if c == nil {
		return
	}

	perQuery := c.searchResults[sp.ID]
	if perQuery == nil {
		perQuery = map[string][]*Task{}
	}
	perQuery[query] = results

	c.searchResults[sp.ID] = perQuery
}

func NewCacheContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, cacheKey, &Cache{
		objs:           map[string]interface{}{},
		spaceIsVisible: map[string]bool{},
		searchResults:  map[string](map[string][]*Task){},
	})
}

func CacheFromContext(ctx context.Context) *Cache {
	c, ok := ctx.Value(cacheKey).(*Cache)
	if !ok {
		return nil
	}

	return c
}

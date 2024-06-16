package namespace

type NamespaceGraphCache struct {
	cache map[string]*NamespaceGraph
}

func NewNamespaceGraphCache() *NamespaceGraphCache {
	cache := new(NamespaceGraphCache)
	cache.cache = make(map[string]*NamespaceGraph)
	return cache
}

func (c *NamespaceGraphCache) Get(key string) (*NamespaceGraph, bool) {
	if val, ok := c.cache[key]; ok {
		return val, true
	}
	return nil, false
}

func (c *NamespaceGraphCache) Put(key string, val *NamespaceGraph) {
	c.cache[key] = val
}

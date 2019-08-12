package two_level_cache

type TwoLevelCache struct {
	ramCache                   *RamCache
	memoryCache                *MemoryCache
	maxRamCapacity             int
	numberOfRequests           int
	numberOfRequestsForRecache int
	structType                 Marshaler
}

func NewTwoLevelCache(maxRamCapacity int, numberOfRequestsForRecache int, marshaler Marshaler) (*TwoLevelCache, error) {
	c := new(TwoLevelCache)
	var err error
	c.memoryCache, err = NewMemoryCache()
	if err != nil {
		return nil, err
	}
	c.ramCache = NewRamCache()
	c.maxRamCapacity = maxRamCapacity
	c.numberOfRequestsForRecache = numberOfRequestsForRecache
	c.structType = marshaler
	return c, nil
}

func (c *TwoLevelCache) Cache(key string, value Marshaler) error {
	if c.ramCache.Size() >= c.maxRamCapacity {
		err := c.memoryCache.Cache(key, value)
		if err != nil {
			return err
		}
	} else {
		err := c.ramCache.Cache(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *TwoLevelCache) Recache() error {
	var meanFrequency int
	for _, k := range append(c.ramCache.GetFrequencies(), c.memoryCache.GetFrequencies()...) {
		meanFrequency += k
	}
	meanFrequency = meanFrequency / c.Size()
	for _, key := range c.ramCache.Keys() {
		freq, _ := c.ramCache.GetFrequency(key)
		if c.ramCache.Size() >= c.maxRamCapacity || freq < meanFrequency {
			obj, _ := c.ramCache.Get(key, c.structType)
			err := c.ramCache.Remove(key)
			if err != nil {
				return err
			}
			err = c.memoryCache.Cache(key, obj)
			if err != nil {
				return err
			}
		}
	}
	for _, key := range c.memoryCache.Keys() {
		freq, _ := c.memoryCache.GetFrequency(key)
		if freq > meanFrequency && c.ramCache.Size() < c.maxRamCapacity {
			obj, _ := c.memoryCache.Get(key, c.structType)
			err := c.memoryCache.Remove(key)
			if err != nil {
				return err
			}
			err = c.ramCache.Cache(key, obj)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *TwoLevelCache) Get(key string) (Marshaler, bool) {
	c.numberOfRequests++
	if c.numberOfRequests >= c.numberOfRequestsForRecache {
		c.numberOfRequests = 0
		err := c.Recache()
		if err != nil {
			return nil, false
		}
	}
	v, ok := c.ramCache.Get(key, c.structType)
	if ok {
		return v, ok
	}
	v, ok = c.memoryCache.Get(key, c.structType)
	if ok {
		return v, ok
	}
	return nil, ok
}

func (c *TwoLevelCache) Remove(key string) error {
	err := c.memoryCache.Remove(key)
	if err != nil {
		return err
	}
	err = c.ramCache.Remove(key)
	if err != nil {
		return err
	}
	return nil
}

func (c *TwoLevelCache) Clear() error {
	err := c.ramCache.Clear()
	if err != nil {
		return err
	}
	err = c.memoryCache.Clear()
	if err != nil {
		return err
	}
	return nil
}

func (c *TwoLevelCache) Contains(key string) bool {
	return c.memoryCache.Contains(key) || c.ramCache.Contains(key)
}

func (c *TwoLevelCache) Size() int {
	return c.ramCache.Size() + c.memoryCache.Size()
}

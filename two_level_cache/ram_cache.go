package two_level_cache

type RamCache struct {
	*objectsMap
	*frequenciesMap
}

func NewRamCache() *RamCache {
	return &RamCache{objectsMap: NewObjectsMap(), frequenciesMap: NewFrequenciesMap()}
}

func (c *RamCache) Cache(key string, val Marshaler) error {
	c.objectsMap.Set(key, val)
	c.SetFrequency(key, 1)
	return nil
}

func (c *RamCache) Get(key string, marshaler Marshaler) (Marshaler, bool) {
	freq, ok := c.GetFrequency(key)
	if !ok {
		return marshaler, false
	}
	c.SetFrequency(key, freq+1)
	var obj Marshaler
	obj, ok = c.objectsMap.Get(key)
	marshaler = obj
	return marshaler, ok
}

func (c *RamCache) Remove(key string) error {
	c.objectsMap.Remove(key)
	c.RemoveFrequency(key)
	return nil
}

func (c *RamCache) Clear() error {
	c.objectsMap.Clear()
	c.frequenciesMap.Clear()
	return nil
}

func (c *RamCache) Contains(key string) bool {
	_, b := c.objectsMap.Get(key)
	return b
}

func (c *RamCache) Size() int {
	return len(c.objectsMap.objects)
}

func (c *RamCache) Keys() []string {
	var keys []string
	for k := range c.objects {
		keys = append(keys, k)
	}
	return keys
}

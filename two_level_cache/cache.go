package two_level_cache

type Cache interface {
	Cache(key string, value Marshaler) error
	Get(key string, marshaller Marshaler) (Marshaler, bool)
	Remove(key string) error
	Clear() error
	Contains(key string) bool
	Size() int
}

package two_level_cache

type Marshaler interface {
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data []byte) error
}

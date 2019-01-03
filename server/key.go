package server

// Key .
type Key struct {
	Name string
	ExpireAt int64
	Type uint8
	Encoding uint8
	value interface{}
}


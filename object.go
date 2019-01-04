package main

// Object .
type Object struct {
	Name string
	ExpireAt int
	Type uint8
	Encoding uint8
	value interface{}
}
// NewObject .
func NewObject() *Object {
	return &Object{

	}
}
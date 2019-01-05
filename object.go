package main

import (
	"time"
)

// Object .
type Object struct {
	Name     string
	ExpireAt int
	Type     uint8
	Encoding uint8
	Value    interface{}
}
// NewObject .
func NewObject() *Object {
	return &Object{
		ExpireAt:-1,
	}
}
// CheckTTL .
func (o *Object) CheckTTL() bool {
	if o.ExpireAt < 0 {
		return true
	}
	t := int(time.Now().UnixNano() / 1e6)
	if t < o.ExpireAt {
		return true
	}

	return false
}
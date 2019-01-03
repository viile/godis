package main

import "sync"

// DB .
type DB struct {
	ID int
	Keys *sync.Map
}
// NewDB .
func NewDB(id int) *DB  {
	return &DB{
		ID:id,
		Keys: &sync.Map{},
	}
}

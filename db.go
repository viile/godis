package main

// DB .
type DB struct {
	ID int
	Keys map[string]*Object
}
// NewDB .
func NewDB(id int) *DB  {
	return &DB{
		ID:id,
		Keys: make(map[string]*Object),
	}
}

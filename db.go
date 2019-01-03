package main

// DB .
type DB struct {
	Id int
	Keys map[string]*Key
}
// NewDB .
func NewDB(id int) *DB  {
	return &DB{
		Id:id,
		Keys: make(map[string]*Key),
	}
}

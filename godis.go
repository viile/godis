package main

import "sync"

// GodisServer .
type GodisServer struct {
	Dbs *sync.Map
}
// NewGodisServer .
func NewGodisServer() *GodisServer {
	s := &GodisServer{
		Dbs: &sync.Map{},
	}
	for i := 0; i < MaxDBNum; i++ {
		s.Dbs.Store(i,NewDB(i))
	}
	return s
}

func (s *GodisServer) Select(db int) (*DB,error){
	r,ok := s.Dbs.Load(db)
	if !ok {
		return nil,DBNotFound
	}
	return r.(*DB),nil
}

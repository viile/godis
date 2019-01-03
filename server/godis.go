package server
// GodisServer .
type GodisServer struct {
	Dbs map[int]*DB
}
// NewGodisServer .
func NewGodisServer() *GodisServer {
	s := &GodisServer{
		Dbs: make(map[int]*DB, MaxDBNum),
	}
	for i := 0; i < MaxDBNum; i++ {
		s.Dbs[i] = NewDB(i)
	}
	return s
}



package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

)

// Server struct
type Server struct {
	// network
	sessions     *sync.Map
	status       int
	listener     net.Listener
	stopCh       chan error
	// object manager
	Dbs *sync.Map
}

// NewServer create a new socket service
func NewServer(addr string) (*Server, error) {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	s := &Server{
		sessions: &sync.Map{},
		stopCh:   make(chan error),
		status:   STInited,
		listener: l,
		Dbs: &sync.Map{},
	}

	for i := 0; i < MaxDBNum; i++ {
		s.Dbs.Store(i,NewDB(i))
	}

	return s, nil
}

// OnMessage HandleMessage .
func (s *Server)OnMessage(ss *Session, buf *[]byte) {
	//fmt.Println("receive msgID:", msg)
	fmt.Println(buf)
	ss.GetConn().SendMessage(nil)
}
// OnDisconnect HandleDisconnect .
func (s *Server) OnDisconnect(ss *Session, err error) {
	fmt.Println(ss.GetConn().GetName() + " lost.")
}
// OnConnect HandleConnect .
func (s *Server) OnConnect(ss *Session) {
	fmt.Println(ss.GetConn().GetName() + " connected.")
	ss.GetConn().SendMessage(nil)
}

// Run Start socket service
func (s *Server) Run() {
	s.status = STRunning
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		s.status = STStop
		cancel()
		s.listener.Close()
	}()

	go s.acceptHandler(ctx)

	for {
		select {
		case <-s.stopCh:
			return
		}
	}
}

func (s *Server) acceptHandler(ctx context.Context) {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			s.stopCh <- err
			return
		}

		go s.connectHandler(ctx, c)
	}
}

func (s *Server) connectHandler(ctx context.Context, c net.Conn) {
	conn := NewConn(c)
	session := NewSession(conn)
	db,err := s.Select(session.settings["db"].(int))
	if err != nil {
		conn.Close()
		return
	}
	session.DBObject = db
	s.sessions.Store(session.GetSessionID(), session)

	connctx, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
		conn.Close()
		s.sessions.Delete(session.GetSessionID())
	}()

	go conn.readCoroutine(connctx)
	go conn.writeCoroutine(connctx)

	s.OnConnect(session)

	for {
		select {
		case err := <-conn.done:
			s.OnDisconnect(session, err)
			return
		case msg := <-conn.messageCh:
			s.OnMessage(session, msg)
		}
	}
}

// GetStatus get socket service status
func (s *Server) GetStatus() int {
	return s.status
}

// Stop stop socket service with reason
func (s *Server) Stop(reason string) {
	s.stopCh <- errors.New(reason)
}

// GetConnsCount get connect count
func (s *Server) GetConnsCount() int {
	var count int
	s.sessions.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count
}

func (s *Server) Select(db int) (*DB,error){
	r,ok := s.Dbs.Load(db)
	if !ok {
		return nil,DBNotFound
	}
	return r.(*DB),nil
}


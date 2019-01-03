package main

import (
	"context"
	"errors"
	"net"
	"sync"

)

// Server struct
type Server struct {
	// handle
	onMessage    func(*Session, *[]byte)
	onConnect    func(*Session)
	onDisconnect func(*Session, error)
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

// RegMessageHandler register message handler
func (s *Server) RegMessageHandler(handler func(*Session, *[]byte)) {
	s.onMessage = handler
}

// RegConnectHandler register connect handler
func (s *Server) RegConnectHandler(handler func(*Session)) {
	s.onConnect = handler
}

// RegDisconnectHandler register disconnect handler
func (s *Server) RegDisconnectHandler(handler func(*Session, error)) {
	s.onDisconnect = handler
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
	s.sessions.Store(session.GetSessionID(), session)

	connctx, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
		conn.Close()
		s.sessions.Delete(session.GetSessionID())
	}()

	go conn.readCoroutine(connctx)
	go conn.writeCoroutine(connctx)

	if s.onConnect != nil {
		s.onConnect(session)
	}

	for {
		select {
		case err := <-conn.done:
			if s.onDisconnect != nil {
				s.onDisconnect(session, err)
			}
			return

		case msg := <-conn.messageCh:
			if s.onMessage != nil {
				s.onMessage(session, msg)
			}
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


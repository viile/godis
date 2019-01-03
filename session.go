package main

import (
	"github.com/google/uuid"
	"log"
)

// Session struct
type Session struct {
	ID       string
	conn     *Conn
	settings map[string]interface{}
	DBObject *DB
	CmdStatus int
}

// NewSession create a new session
func NewSession(conn *Conn) *Session {
	session := &Session{
		ID:       uuid.New().String(),
		conn:     conn,
		settings: make(map[string]interface{}),
	}
	session.settings["db"] = 0
	session.settings["auth"] = true
	return session
}

// GetSessionID get session ID
func (s *Session) GetSessionID() string {
	return s.ID
}

// OnDisconnect .
func (s *Session) OnDisconnect(err error) {
	log.Println(s.conn.GetName() + " lost.\n",err)
}
// OnConnect .
func (s *Session) OnConnect() {
	log.Println(s.conn.GetName() + " connected.")
	pkg := []byte{43, 79, 75, 13, 10}
	s.conn.SendMessage(pkg)
}

// OnHandle .
func (s *Session) OnHandle(buf *[]byte) {
	log.Println(buf)
	pkg := []byte{43, 79, 75, 13, 10}
	s.conn.SendMessage(pkg)
}



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
	Parser *Codec
}

// NewSession create a new session
func NewSession(conn *Conn) *Session {
	session := &Session{
		ID:       uuid.New().String(),
		conn:     conn,
		settings: make(map[string]interface{}),
		Parser:NewCodec(),
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
	log.Println(s.conn.GetName() + " lost")
}
// OnConnect .
func (s *Session) OnConnect() {
	log.Println(s.conn.GetName() + " connected.")
}

// OnHandle .
func (s *Session) OnHandle(buf *[]byte) {
	resps := s.Parser.Decode(buf)
	for _,v := range resps{
		if v.Argc == 0  {
			continue
		}
		ret := HM.Distribute(s.DBObject,v)
		if len(ret) == 0 {
			continue
		}
		s.conn.SendMessage(ret)
	}
}


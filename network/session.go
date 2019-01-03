package network

import (
	uuid "github.com/google/uuid"
)

// Session struct
type Session struct {
	sID      string
	conn     *Conn
	settings map[string]interface{}
}

// NewSession create a new session
func NewSession(conn *Conn) *Session {
	id := uuid.New()
	session := &Session{
		sID:      id.String(),
		conn:     conn,
		settings: make(map[string]interface{}),
	}
	session.settings["db"] = 0
	return session
}

// GetSessionID get session ID
func (s *Session) GetSessionID() string {
	return s.sID
}

// GetConn get zero.Conn pointer
func (s *Session) GetConn() *Conn {
	return s.conn
}

// SetConn set a zero.Conn to session
func (s *Session) SetConn(conn *Conn) {
	s.conn = conn
}

// GetSetting get setting
func (s *Session) GetSetting(key string) interface{} {

	if v, ok := s.settings[key]; ok {
		return v
	}

	return nil
}

// SetSetting set setting
func (s *Session) SetSetting(key string, value interface{}) {
	s.settings[key] = value
}

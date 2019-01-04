package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"strings"
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
	log.Println(s.conn.GetName() + " lost.\n")
}
// OnConnect .
func (s *Session) OnConnect() {
	log.Println(s.conn.GetName() + " connected.")
}

// OnHandle .
func (s *Session) OnHandle(buf *[]byte) {
	resps := s.Parser.Decode(buf)
	for _,v := range resps{
		if v.Argc < 1  {
			continue
		}
		ret := s.Handle(v)
		if len(ret) == 0 {
			continue
		}
		s.conn.SendMessage(ret)
	}
}

func (s *Session) Handle(resp *Resp) (ret []byte) {
	ret = make([]byte,0)
	log.Println(resp)
	cmd := strings.ToUpper(resp.Argv[0])
	switch cmd {
	case DEL:
		if resp.Argc < 2 {
			ret =s.Parser.ErrReplyEncode(fmt.Sprintf(ErrCommandArgsWrongNumber.Error(),cmd))
			return
		}
		ret = s.Parser.IntReplyEncode(s.DBObject.Del(resp.Argv[1:]))
		return
	case SET:
		if resp.Argc < 3 {
			ret =s.Parser.ErrReplyEncode(fmt.Sprintf(ErrCommandArgsWrongNumber.Error(),cmd))
			return
		}
		err := s.DBObject.Set(resp.Argv[1:])
		if err == nil {
			ret = s.Parser.StatusOkReplyEncode()
		}else {
			ret = s.Parser.NilBulkReplyEncode()
		}
		return
	case EXPIRE:
		if resp.Argc < 3 {
			ret =s.Parser.ErrReplyEncode(fmt.Sprintf(ErrCommandArgsWrongNumber.Error(),cmd))
			return
		}
		ret = s.Parser.IntReplyEncode(s.DBObject.EXPIRE(resp.Argv[1:]))
		return
	case TTL:
		if resp.Argc < 2 {
			ret =s.Parser.ErrReplyEncode(fmt.Sprintf(ErrCommandArgsWrongNumber.Error(),cmd))
			return
		}
		ret = s.Parser.IntReplyEncode(s.DBObject.TTL(resp.Argv[1]))
		return
	case PTTL:
		if resp.Argc < 2 {
			ret =s.Parser.ErrReplyEncode(fmt.Sprintf(ErrCommandArgsWrongNumber.Error(),cmd))
			return
		}
		ret = s.Parser.IntReplyEncode(s.DBObject.PTTL(resp.Argv[1]))
		return
	case GET:
		if resp.Argc < 2 {
			ret = s.Parser.ErrReplyEncode(fmt.Sprintf(ErrCommandArgsWrongNumber.Error(),cmd))
			return
		}
		r,err := s.DBObject.Get(resp.Argv[1])
		log.Println(r,err)
		if err == nil {
			ret = s.Parser.BulkReplyEncode(r)
		} else if err == ErrTypeNotMatch {
			ret = s.Parser.ErrReplyEncode(err.Error())
		} else if err == ErrKeyNotFound {
			ret = s.Parser.NilBulkReplyEncode()
		}
		return
	default:
		 ret = s.Parser.ErrReplyEncode(ErrDontSupportThisCommand.Error())
	}
	return
}


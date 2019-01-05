package main

import (
	"fmt"
	"strings"
)

var HM *HandleManger
type Handle struct {
	Command string
	LeastArgNumbers int
	Func func(*DB,*Resp) []byte
}

type HandleManger struct {
	Handles map[string]*Handle
}

func init(){
	HM = NewHandleManger()
	HM.Handles[GET] = &Handle{
		Command:GET,
		LeastArgNumbers:2,
		Func:Get,
	}
	HM.Handles[SET] = &Handle{
		Command:SET,
		LeastArgNumbers:3,
		Func:Set,
	}
	HM.Handles[DEL] = &Handle{
		Command:DEL,
		LeastArgNumbers:2,
		Func:Del,
	}
	HM.Handles[TTL] = &Handle{
		Command:TTL,
		LeastArgNumbers:2,
		Func:Ttl,
	}
	HM.Handles[PTTL] = &Handle{
		Command:PTTL,
		LeastArgNumbers:2,
		Func:PTtl,
	}
	HM.Handles[EXISTS] = &Handle{
		Command:EXISTS,
		LeastArgNumbers:2,
		Func:Exists,
	}
	HM.Handles[EXPIRE] = &Handle{
		Command:EXPIRE,
		LeastArgNumbers:3,
		Func:Expire,
	}
	HM.Handles[PEXPIRE] = &Handle{
		Command:PEXPIRE,
		LeastArgNumbers:3,
		Func:PExpire,
	}
	HM.Handles[EXPIREAT] = &Handle{
		Command:EXPIREAT,
		LeastArgNumbers:3,
		Func:ExpireAt,
	}
	HM.Handles[PEXPIREAT] = &Handle{
		Command:PEXPIREAT,
		LeastArgNumbers:3,
		Func:PExpireAt,
	}
	HM.Handles[PERSIST] = &Handle{
		Command:PERSIST,
		LeastArgNumbers:3,
		Func:Persist,
	}
}

func NewHandleManger() *HandleManger{
	return &HandleManger{
		Handles:make(map[string]*Handle),
	}
}

func (m *HandleManger) Distribute(db *DB,resp *Resp) []byte {
	cmd := strings.ToUpper(resp.Argv[0])
	handle,ok := m.Handles[cmd]
	if !ok {
		return ErrReplyEncode(ErrDontSupportThisCommand.Error())
	}
	if resp.Argc < handle.LeastArgNumbers {
		return ErrReplyEncode(fmt.Sprintf(ErrCommandArgsWrongNumber.Error(),cmd))
	}
	return handle.Func(db,resp)
}


package main

import (
	"fmt"
	"strings"
)

var HM *HandleManger

type Handle struct {
	Command       string
	MinArgNumbers int
	MaxArgNumbers int
	Func          func(*DB, *Resp) []byte
}

type HandleManger struct {
	Handles map[string]*Handle
}

func init() {
	HM = NewHandleManger()
	HM.Handles[GET] = &Handle{
		Command:       GET,
		MinArgNumbers: 2,
		Func:          Get,
	}
	HM.Handles[SET] = &Handle{
		Command:       SET,
		MinArgNumbers: 3,
		Func:          Set,
	}
	HM.Handles[DEL] = &Handle{
		Command:       DEL,
		MinArgNumbers: 2,
		Func:          Del,
	}
	HM.Handles[TTL] = &Handle{
		Command:       TTL,
		MinArgNumbers: 2,
		Func:          Ttl,
	}
	HM.Handles[PTTL] = &Handle{
		Command:       PTTL,
		MinArgNumbers: 2,
		Func:          PTtl,
	}
	HM.Handles[EXISTS] = &Handle{
		Command:       EXISTS,
		MinArgNumbers: 2,
		Func:          Exists,
	}
	HM.Handles[EXPIRE] = &Handle{
		Command:       EXPIRE,
		MinArgNumbers: 3,
		Func:          Expire,
	}
	HM.Handles[PEXPIRE] = &Handle{
		Command:       PEXPIRE,
		MinArgNumbers: 3,
		Func:          PExpire,
	}
	HM.Handles[EXPIREAT] = &Handle{
		Command:       EXPIREAT,
		MinArgNumbers: 3,
		Func:          ExpireAt,
	}
	HM.Handles[PEXPIREAT] = &Handle{
		Command:       PEXPIREAT,
		MinArgNumbers: 3,
		Func:          PExpireAt,
	}
	HM.Handles[PERSIST] = &Handle{
		Command:       PERSIST,
		MinArgNumbers: 3,
		Func:          Persist,
	}
	HM.Handles[PING] = &Handle{
		Command:       PING,
		MinArgNumbers: 1,
		MaxArgNumbers: 2,
		Func:          Ping,
	}
	HM.Handles[ECHO] = &Handle{
		Command:       ECHO,
		MinArgNumbers: 2,
		MaxArgNumbers: 2,
		Func:          Echo,
	}
}

func NewHandleManger() *HandleManger {
	return &HandleManger{
		Handles: make(map[string]*Handle),
	}
}

func (m *HandleManger) Distribute(db *DB, resp *Resp) []byte {
	cmd := strings.ToUpper(resp.Argv[0])
	handle, ok := m.Handles[cmd]
	if !ok {
		return ErrReplyEncode(ErrDontSupportThisCommand.Error())
	}
	if resp.Argc < handle.MinArgNumbers || (handle.MaxArgNumbers > 0 && resp.Argc > handle.MaxArgNumbers) {
		return ErrReplyEncode(fmt.Sprintf(ErrCommandArgsWrongNumber.Error(), cmd))
	}
	return handle.Func(db, resp)
}
package main

import (
	"fmt"
	"strings"
)
// HM .
var HM *HandleManger
// Handle .
type Handle struct {
	Command       string
	MinArgNumbers int
	MaxArgNumbers int
	Func          func(*DB, *Resp) []byte
}
// HandleManger .
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
	HM.Handles[HSET] = &Handle{
		Command:       HSET,
		MinArgNumbers: 4,
		MaxArgNumbers: 4,
		Func:          HSet,
	}
	HM.Handles[HGET] = &Handle{
		Command:       HGET,
		MinArgNumbers: 3,
		MaxArgNumbers: 3,
		Func:          HGet,
	}
	HM.Handles[HGETALL] = &Handle{
		Command:       HGETALL,
		MinArgNumbers: 2,
		MaxArgNumbers: 2,
		Func:          HGetAll,
	}
	HM.Handles[HDEL] = &Handle{
		Command:       HDEL,
		MinArgNumbers: 3,
		Func:          HDel,
	}
	HM.Handles[HEXISTS] = &Handle{
		Command:       HEXISTS,
		MinArgNumbers: 3,
		MaxArgNumbers: 3,
		Func:          HExists,
	}
	HM.Handles[HLEN] = &Handle{
		Command:       HLEN,
		MinArgNumbers: 2,
		MaxArgNumbers: 2,
		Func:          HLen,
	}
	HM.Handles[HKEYS] = &Handle{
		Command:       HKEYS,
		MinArgNumbers: 2,
		MaxArgNumbers: 2,
		Func:          HKeys,
	}
	HM.Handles[HVALS] = &Handle{
		Command:       HVALS,
		MinArgNumbers: 2,
		MaxArgNumbers: 2,
		Func:          HVals,
	}
	HM.Handles[HINCRBY] = &Handle{
		Command:       HINCRBY,
		MinArgNumbers: 4,
		MaxArgNumbers: 4,
		Func:          HIncrBy,
	}
	HM.Handles[HINCRBYFLOAT] = &Handle{
		Command:       HINCRBYFLOAT,
		MinArgNumbers: 4,
		MaxArgNumbers: 4,
		Func:          HIncrByFloat,
	}
	HM.Handles[HSETNX] = &Handle{
		Command:       HSETNX,
		MinArgNumbers: 4,
		MaxArgNumbers: 4,
		Func:          HSetNX,
	}
	HM.Handles[SADD] = &Handle{
		Command:       SADD,
		MinArgNumbers: 3,
		Func:          SAdd,
	}
	HM.Handles[SCARD] = &Handle{
		Command:       SCARD,
		MinArgNumbers: 2,
		MaxArgNumbers: 2,
		Func:          SCard,
	}
	HM.Handles[SDIFF] = &Handle{
		Command:       SDIFF,
		MinArgNumbers: 2,
		Func:          SDiff,
	}
	HM.Handles[SDIFFSTORE] = &Handle{
		Command:       SDIFFSTORE,
		MinArgNumbers: 3,
		Func:          SDiffStore,
	}
	HM.Handles[SINTER] = &Handle{
		Command:       SINTER,
		MinArgNumbers: 2,
		Func:          SInter,
	}
	HM.Handles[SINTERSTORE] = &Handle{
		Command:       SINTERSTORE,
		MinArgNumbers: 3,
		Func:          SInterStore,
	}
	HM.Handles[SISMEMBER] = &Handle{
		Command:       SISMEMBER,
		MinArgNumbers: 3,
		MaxArgNumbers: 3,
		Func:          SIsMember,
	}
	HM.Handles[SMEMBERS] = &Handle{
		Command:       SMEMBERS,
		MinArgNumbers: 2,
		MaxArgNumbers: 2,
		Func:          SMembers,
	}
	HM.Handles[SMOVE] = &Handle{
		Command:       SMOVE,
		MinArgNumbers: 4,
		MaxArgNumbers: 4,
		Func:          SMove,
	}
	HM.Handles[SPOP] = &Handle{
		Command:       SPOP,
		MinArgNumbers: 2,
		MaxArgNumbers: 2,
		Func:          SPop,
	}
	HM.Handles[SREM] = &Handle{
		Command:       SREM,
		MinArgNumbers: 3,
		Func:          SRem,
	}
	HM.Handles[SUNION] = &Handle{
		Command:       SUNION,
		MinArgNumbers: 2,
		Func:          SUnion,
	}
	HM.Handles[SUNIONSTORE] = &Handle{
		Command:       SUNIONSTORE,
		MinArgNumbers: 3,
		Func:          SUnionStore,
	}
}
// NewHandleManger .
func NewHandleManger() *HandleManger {
	return &HandleManger{
		Handles: make(map[string]*Handle),
	}
}
// Distribute .
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

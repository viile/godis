package main

import (
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DB .
type DB struct {
	ID      int
	Objects *sync.Map
}
// NewDB .
func NewDB(id int) *DB  {
	return &DB{
		ID:      id,
		Objects: &sync.Map{},
	}
}

func(d *DB) GetObject(key string) (*Object,error){
	v,ok := d.Objects.Load(key)
	if !ok {
		return nil,ErrKeyNotFound
	}
	object := v.(*Object)
	if !object.CheckTTL() {
		d.Objects.Delete(key)
		return nil,ErrKeyNotFound
	}
	return object,nil
}
func Exists(d *DB,resp *Resp) []byte {
	_,err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}
	return IntReplyEncode(1)
}
func Del(d *DB,resp *Resp) []byte {
	ret := 0
	for _,v := range resp.Argv[1:] {
		_,ok := d.Objects.Load(v)
		if !ok {
			continue
		}
		d.Objects.Delete(v)
		ret++
	}
	return IntReplyEncode(ret)
}
func Set(d *DB,resp *Resp) []byte {
	var NXFlag,XXFlag bool
	EXValue := -1
	key,value := resp.Argv[1],resp.Argv[2]
	if resp.Argc > 3 {
		index := 3
		ctime := int(time.Now().UnixNano() / 1e6)
		for index < resp.Argc {
			switch strings.ToUpper(resp.Argv[index]) {
			case "NX":
				NXFlag = true
				index++
			case "XX":
				XXFlag = true
				index++
			case "PX":
				if index + 1 >= resp.Argc {
					return ErrReplyEncode(ErrCommand.Error())
				}
				t := resp.Argv[index + 1]
				r,err:= strconv.Atoi(t)
				if err != nil{
					return ErrReplyEncode(ErrCommand.Error())
				}
				EXValue = ctime + r
				index+= 2
			case "EX":
				if index + 1 >= resp.Argc {
					return ErrReplyEncode(ErrCommand.Error())
				}
				t := resp.Argv[index + 1]
				r,err:= strconv.Atoi(t)
				if err != nil{
					return ErrReplyEncode(ErrCommand.Error())
				}
				EXValue = ctime + r * 1000
				index+= 2
			default:
				return ErrReplyEncode(ErrCommand.Error())
			}
		}
	}

	_,ok := d.Objects.Load(key)
	if (ok && NXFlag) || (!ok && XXFlag) {
		return NilBulkReplyEncode()
	}

	o := NewObject()
	o.Value = NewRedisString(value)
	o.Type = TypeRedisString
	o.Encoding = RedisEncodingRaw
	o.Name = key
	o.ExpireAt = EXValue

	d.Objects.Store(key,o)

	return StatusOkReplyEncode()
}
func Get(d *DB,resp *Resp) []byte {
	object,err := d.GetObject(resp.GetKey())
	if err != nil {
		return ErrReplyEncode(ErrKeyNotFound.Error())
	}
	if object.Type != TypeRedisString {
		return ErrReplyEncode(ErrTypeNotMatch.Error())
	}
	v := object.Value.(*RedisString)
	return BulkReplyEncode(v.value)
}
func Type(d *DB,resp *Resp) []byte {
	ret := ""
	object,err := d.GetObject(resp.GetKey())
	if err != nil {
		ret = "none"
	}
	switch object.Type {
	case TypeRedisString:
		ret = "string"
	case TypeRedisHash:
		ret = "hash"
	case TypeRedisList:
		ret = "list"
	case TypeRedisSet:
		ret = "set"
	case TypeRedisZSet:
		ret = "zset"
	default:
		ret = "none"
	}

	return BulkReplyEncode(ret)
}
func ExpireAt(d *DB,resp *Resp) []byte {
	object,err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}

	e := resp.Argv[2]
	ex,err := strconv.Atoi(e)
	if err != nil {
		return IntReplyEncode(0)
	}

	object.ExpireAt = ex * 1000

	return IntReplyEncode(1)
}
func PExpireAt(d *DB,resp *Resp) []byte {
	object,err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}

	e := resp.Argv[2]
	ex,err := strconv.Atoi(e)
	if err != nil {
		return IntReplyEncode(0)
	}

	object.ExpireAt = ex

	return IntReplyEncode(1)
}
func Expire(d *DB,resp *Resp) []byte {
	object,err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}
	e := resp.Argv[2]
	ex,err := strconv.Atoi(e)
	if err != nil {
		return IntReplyEncode(0)
	}

	object.ExpireAt = int(time.Now().UnixNano() / 1e6) + ex * 1000

	return IntReplyEncode(1)
}
func PExpire(d *DB,resp *Resp) []byte {
	object,err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}

	ex,err := strconv.Atoi(resp.Argv[2])
	if err != nil {
		return IntReplyEncode(0)
	}

	object.ExpireAt = int(time.Now().UnixNano() / 1e6) + ex

	return IntReplyEncode(1)
}

func Persist(d *DB,resp *Resp) []byte {
	object,err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}

	object.ExpireAt = -1

	return IntReplyEncode(1)
}
func Ttl(d *DB,resp *Resp) []byte {
	object,err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(KeyNotExists)
	}

	if object.ExpireAt < 0 {
		return IntReplyEncode(object.ExpireAt)
	}else {
		t := object.ExpireAt - int(time.Now().UnixNano() / 1e6)
		return IntReplyEncode(int(math.Ceil(float64(t) / 1000.0)))
	}
}
func PTtl(d *DB,resp *Resp) []byte {
	object,err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(KeyNotExists)
	}

	if object.ExpireAt < 0 {
		return IntReplyEncode(object.ExpireAt)
	}else {
		t := object.ExpireAt - int(time.Now().UnixNano() / 1e6)
		return IntReplyEncode(t)
	}
}

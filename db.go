package main

import (
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// DB .
type DB struct {
	ID      int
	Objects *sync.Map
	Length int32
	Server  *Server
}

// NewDB .
func NewDB(id int, server *Server) *DB {
	return &DB{
		ID:      id,
		Server:  server,
		Objects: &sync.Map{},
	}
}

func (d *DB) AddObject(key string,obj *Object) (error) {
	d.Objects.Store(key,obj)
	//d.Length++
	atomic.AddInt32(&d.Length,1)
	return nil
}
func (d *DB) DelObject(key string) (error) {
	d.Objects.Delete(key)
	atomic.AddInt32(&d.Length,-1)
	//d.Length--
	return nil
}

func (d *DB) Flush() (error) {
	d.Objects = &sync.Map{}
	d.Length = 0
	return nil
}

func (d *DB) GetObject(key string) (*Object, error) {
	v, ok := d.Objects.Load(key)
	if !ok {
		return nil, ErrKeyNotFound
	}
	object := v.(*Object)
	if !object.CheckTTL() {
		d.DelObject(key)
		return nil, ErrKeyNotFound
	}
	return object, nil
}
func Exists(d *DB, resp *Resp) []byte {
	_, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}
	return IntReplyEncode(1)
}
func Del(d *DB, resp *Resp) []byte {
	ret := 0
	for _, v := range resp.Argv[1:] {
		_, err := d.GetObject(v)
		if err != nil {
			continue
		}
		d.DelObject(v)
		ret++
	}
	return IntReplyEncode(ret)
}
func Set(d *DB, resp *Resp) []byte {
	var NXFlag, XXFlag bool
	EXValue := -1
	key, value := resp.Argv[1], resp.Argv[2]
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
				if index+1 >= resp.Argc {
					return ErrReplyEncode(ErrCommand.Error())
				}
				t := resp.Argv[index+1]
				r, err := strconv.Atoi(t)
				if err != nil {
					return ErrReplyEncode(ErrCommand.Error())
				}
				EXValue = ctime + r
				index += 2
			case "EX":
				if index+1 >= resp.Argc {
					return ErrReplyEncode(ErrCommand.Error())
				}
				t := resp.Argv[index+1]
				r, err := strconv.Atoi(t)
				if err != nil {
					return ErrReplyEncode(ErrCommand.Error())
				}
				EXValue = ctime + r*1000
				index += 2
			default:
				return ErrReplyEncode(ErrCommand.Error())
			}
		}
	}

	_, err := d.GetObject(key)
	if (err == nil && NXFlag) || (err != nil && XXFlag) {
		return NilBulkReplyEncode()
	}

	o := NewObject()
	o.Value = NewRedisString(value)
	o.Type = TypeRedisString
	o.Encoding = RedisEncodingRaw
	o.Name = key
	o.ExpireAt = EXValue

	d.AddObject(key, o)

	return StatusOkReplyEncode()
}
func Get(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return ErrReplyEncode(ErrKeyNotFound.Error())
	}
	if object.Type != TypeRedisString {
		return ErrReplyEncode(ErrTypeNotMatch.Error())
	}
	v := object.Value.(*RedisString)
	return BulkReplyEncode(v.value)
}
func Type(d *DB, resp *Resp) []byte {
	ret := ""
	object, err := d.GetObject(resp.GetKey())
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
func ExpireAt(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}

	e := resp.Argv[2]
	ex, err := strconv.Atoi(e)
	if err != nil {
		return IntReplyEncode(0)
	}

	object.ExpireAt = ex * 1000

	return IntReplyEncode(1)
}
func PExpireAt(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}

	e := resp.Argv[2]
	ex, err := strconv.Atoi(e)
	if err != nil {
		return IntReplyEncode(0)
	}

	object.ExpireAt = ex

	return IntReplyEncode(1)
}
func Expire(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}
	e := resp.Argv[2]
	ex, err := strconv.Atoi(e)
	if err != nil {
		return IntReplyEncode(0)
	}

	object.ExpireAt = int(time.Now().UnixNano()/1e6) + ex*1000

	return IntReplyEncode(1)
}
func PExpire(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}

	ex, err := strconv.Atoi(resp.Argv[2])
	if err != nil {
		return IntReplyEncode(0)
	}

	object.ExpireAt = int(time.Now().UnixNano()/1e6) + ex

	return IntReplyEncode(1)
}

func Persist(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}

	object.ExpireAt = -1

	return IntReplyEncode(1)
}
func Ttl(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(KeyNotExists)
	}

	if object.ExpireAt < 0 {
		return IntReplyEncode(object.ExpireAt)
	} else {
		t := object.ExpireAt - int(time.Now().UnixNano()/1e6)
		return IntReplyEncode(int(math.Ceil(float64(t) / 1000.0)))
	}
}
func PTtl(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(KeyNotExists)
	}

	if object.ExpireAt < 0 {
		return IntReplyEncode(object.ExpireAt)
	} else {
		t := object.ExpireAt - int(time.Now().UnixNano()/1e6)
		return IntReplyEncode(t)
	}
}
func Ping(d *DB, resp *Resp) []byte {
	if resp.Argc == 2 {
		return BulkReplyEncode(resp.Argv[1])
	}
	return BulkReplyEncode("PONG")
}
func Echo(d *DB, resp *Resp) []byte {
	return BulkReplyEncode(resp.Argv[1])
}

func HSet(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	var hash *RedisHash
	if err != nil {
		hash = NewRedisHash()
		object = NewObject()
		object.Type = TypeRedisHash
		object.Name = resp.GetKey()
		object.Encoding = RedisEncodingHt
		object.Value = hash
		d.AddObject(resp.GetKey(), object)
	} else {
		if object.Type != TypeRedisHash {
			return ErrReplyEncode(ErrWrongKeyType.Error())
		}
		hash = object.Value.(*RedisHash)
	}

	ret := hash.Set(resp.Argv[2], resp.Argv[3])
	return IntReplyEncode(ret)
}
func HSetNX(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	var hash *RedisHash
	if err != nil {
		hash = NewRedisHash()
		object = NewObject()
		object.Type = TypeRedisHash
		object.Name = resp.GetKey()
		object.Encoding = RedisEncodingHt
		object.Value = hash
		d.AddObject(resp.GetKey(), object)
	} else {
		if object.Type != TypeRedisHash {
			return ErrReplyEncode(ErrWrongKeyType.Error())
		}
		hash = object.Value.(*RedisHash)
	}

	ret := hash.SetNX(resp.Argv[2], resp.Argv[3])
	return IntReplyEncode(ret)
}

func HGet(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisHash {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	hash := object.Value.(*RedisHash)
	ret, err := hash.Get(resp.Argv[2])
	if err != nil {
		return NilBulkReplyEncode()
	}
	return BulkReplyEncode(ret)
}

func HDel(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisHash {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	obj := object.Value.(*RedisHash)
	ret := obj.Del(resp.Argv[2:])
	return IntReplyEncode(ret)
}

func HExists(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisHash {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	obj := object.Value.(*RedisHash)
	ret := obj.Exists(resp.Argv[2])
	return IntReplyEncode(ret)
}
func HLen(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisHash {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	obj := object.Value.(*RedisHash)
	ret := obj.Len()
	return IntReplyEncode(ret)
}

func HGetAll(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisHash {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	obj := object.Value.(*RedisHash)
	ret := obj.GetAll()
	return ArrayReplyEncode(ret)
}
func HKeys(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisHash {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	obj := object.Value.(*RedisHash)
	ret := obj.Keys()
	return ArrayReplyEncode(ret)
}
func HVals(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisHash {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	obj := object.Value.(*RedisHash)
	ret := obj.Vals()
	return ArrayReplyEncode(ret)
}

func HIncrBy(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisHash {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	value, err := strconv.Atoi(resp.Argv[3])
	if err != nil {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	hash := object.Value.(*RedisHash)
	ret, err := hash.IncrBy(resp.Argv[2], value)
	if err != nil {
		return NilBulkReplyEncode()
	}
	return IntReplyEncode(ret)
}
func HIncrByFloat(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisHash {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	value, err := strconv.ParseFloat(resp.Argv[3], 64)
	if err != nil {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	hash := object.Value.(*RedisHash)
	ret, err := hash.IncrByFloat(resp.Argv[2], value)
	if err != nil {
		return NilBulkReplyEncode()
	}
	return BulkReplyEncode(strconv.FormatFloat(ret, 'f', -1, 64))
}

func SAdd(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	var set *RedisSet
	if err != nil {
		set = NewRedisSet()
		object = NewObject()
		object.Type = TypeRedisSet
		object.Name = resp.GetKey()
		object.Encoding = RedisEncodingHt
		object.Value = set
		d.AddObject(resp.GetKey(), object)
	} else {
		if object.Type != TypeRedisSet {
			return ErrReplyEncode(ErrWrongKeyType.Error())
		}
	}

	ret := set.Add(resp.Argv[2:])
	return IntReplyEncode(ret)
}

func SCard(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}
	if object.Type != TypeRedisSet {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	set := object.Value.(*RedisSet)
	return IntReplyEncode(set.Card())
}

func SDiff(d *DB, resp *Resp) []byte {
	var ret *RedisSet
	for _, v := range resp.Argv[1:] {
		object, err := d.GetObject(v)
		if err != nil || object.Type != TypeRedisSet {
			continue
		}
		if ret == nil {
			ret = object.Value.(*RedisSet)
		} else {
			ret = ret.Diff(object.Value.(*RedisSet))
		}
	}
	if ret == nil {
		return NilBulkReplyEncode()
	} else {
		return ArrayReplyEncode(ret.Members())
	}
}

func SDiffStore(d *DB, resp *Resp) []byte {
	var ret *RedisSet
	for _, v := range resp.Argv[2:] {
		object, err := d.GetObject(v)
		if err != nil || object.Type != TypeRedisSet {
			continue
		}
		if ret == nil {
			ret = object.Value.(*RedisSet)
		} else {
			ret = ret.Diff(object.Value.(*RedisSet))
		}
	}
	if ret == nil {
		ret = NewRedisSet()
	}
	object := NewObject()
	object.Type = TypeRedisSet
	object.Name = resp.GetKey()
	object.Encoding = RedisEncodingHt
	object.Value = ret
	d.AddObject(resp.GetKey(), object)
	return IntReplyEncode(ret.Card())
}

func SInter(d *DB, resp *Resp) []byte {
	var ret *RedisSet
	for _, v := range resp.Argv[1:] {
		object, err := d.GetObject(v)
		if err != nil || object.Type != TypeRedisSet {
			continue
		}
		if ret == nil {
			ret = object.Value.(*RedisSet)
		} else {
			ret = ret.Inter(object.Value.(*RedisSet))
		}
	}
	if ret == nil {
		return NilBulkReplyEncode()
	} else {
		return ArrayReplyEncode(ret.Members())
	}
}

func SInterStore(d *DB, resp *Resp) []byte {
	var ret *RedisSet
	for _, v := range resp.Argv[2:] {
		object, err := d.GetObject(v)
		if err != nil || object.Type != TypeRedisSet {
			continue
		}
		if ret == nil {
			ret = object.Value.(*RedisSet)
		} else {
			ret = ret.Inter(object.Value.(*RedisSet))
		}
	}
	if ret == nil {
		ret = NewRedisSet()
	}
	object := NewObject()
	object.Type = TypeRedisSet
	object.Name = resp.GetKey()
	object.Encoding = RedisEncodingHt
	object.Value = ret
	d.AddObject(resp.GetKey(), object)
	return IntReplyEncode(ret.Card())
}

func SIsMember(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}

	if object.Type != TypeRedisSet {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	set := object.Value.(*RedisSet)
	if !set.IsMember(resp.Argv[2]) {
		return IntReplyEncode(0)
	}
	return IntReplyEncode(1)
}

func SMembers(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}

	if object.Type != TypeRedisSet {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	set := object.Value.(*RedisSet)
	return ArrayReplyEncode(set.Members())
}

func SMove(d *DB, resp *Resp) []byte {
	src, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}
	if src.Type != TypeRedisSet {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	srcset := src.Value.(*RedisSet)
	if !srcset.IsMember(resp.Argv[3]) {
		return IntReplyEncode(0)
	}
	dst, err := d.GetObject(resp.Argv[2])
	if err != nil {
		return IntReplyEncode(0)
	}
	if dst.Type != TypeRedisSet {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	dstset := dst.Value.(*RedisSet)
	srcset.Move(dstset, resp.Argv[3])
	return IntReplyEncode(1)
}

func SPop(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}

	if object.Type != TypeRedisSet {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	set := object.Value.(*RedisSet)
	return BulkReplyEncode(set.Pop())
}

func SRem(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}

	if object.Type != TypeRedisSet {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	set := object.Value.(*RedisSet)
	return IntReplyEncode(set.Rem(resp.Argv[2:]))
}

func SUnion(d *DB, resp *Resp) []byte {
	var ret *RedisSet
	for _, v := range resp.Argv[1:] {
		object, err := d.GetObject(v)
		if err != nil || object.Type != TypeRedisSet {
			continue
		}
		if ret == nil {
			ret = object.Value.(*RedisSet)
		} else {
			ret = ret.Union(object.Value.(*RedisSet))
		}
	}
	if ret == nil {
		return NilBulkReplyEncode()
	} else {
		return ArrayReplyEncode(ret.Members())
	}
}

func SUnionStore(d *DB, resp *Resp) []byte {
	var ret *RedisSet
	for _, v := range resp.Argv[2:] {
		object, err := d.GetObject(v)
		if err != nil || object.Type != TypeRedisSet {
			continue
		}
		if ret == nil {
			ret = object.Value.(*RedisSet)
		} else {
			ret = ret.Union(object.Value.(*RedisSet))
		}
	}
	if ret == nil {
		ret = NewRedisSet()
	}
	object := NewObject()
	object.Type = TypeRedisSet
	object.Name = resp.GetKey()
	object.Encoding = RedisEncodingHt
	object.Value = ret
	d.AddObject(resp.GetKey(), object)
	return IntReplyEncode(ret.Card())
}

func LPush(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	var l *RedisList
	if err != nil {
		l = NewRedisList()
		object := NewObject()
		object.Type = TypeRedisList
		object.Name = resp.GetKey()
		object.Encoding = RedisEncodingLinkedList
		object.Value = l
		d.AddObject(resp.GetKey(), object)
	} else {
		if object.Type != TypeRedisList {
			return ErrReplyEncode(ErrWrongKeyType.Error())
		}
		l = object.Value.(*RedisList)
	}

	return IntReplyEncode(l.Push(resp.Argv[2:]))
}
func LPushx(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	var l *RedisList
	if err != nil {
		return IntReplyEncode(0)
	}
	if object.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	l = object.Value.(*RedisList)
	if l.Len() == 0 {
		return IntReplyEncode(0)
	}

	return IntReplyEncode(l.Push(resp.Argv[2:]))
}
func RPush(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	var l *RedisList
	if err != nil {
		l = NewRedisList()
		object := NewObject()
		object.Type = TypeRedisList
		object.Name = resp.GetKey()
		object.Encoding = RedisEncodingLinkedList
		object.Value = l
		d.AddObject(resp.GetKey(), object)
	} else {
		if object.Type != TypeRedisList {
			return ErrReplyEncode(ErrWrongKeyType.Error())
		}
		l = object.Value.(*RedisList)
	}

	return IntReplyEncode(l.RPush(resp.Argv[2:]))
}

func RPushHx(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	var l *RedisList
	if err != nil {
		return IntReplyEncode(0)
	}
	if object.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	l = object.Value.(*RedisList)
	if l.Len() == 0 {
		return IntReplyEncode(0)
	}

	return IntReplyEncode(l.RPush(resp.Argv[2:]))
}

func LPop(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	l := object.Value.(*RedisList)
	if l.Len() == 0 {
		return NilBulkReplyEncode()
	}
	return BulkReplyEncode(l.LPop())
}

func RPop(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	l := object.Value.(*RedisList)
	if l.Len() == 0 {
		return NilBulkReplyEncode()
	}
	return BulkReplyEncode(l.RPop())
}

func RPopLPush(d *DB, resp *Resp) []byte {
	src, err := d.GetObject(resp.Argv[1])
	if err != nil {
		return NilBulkReplyEncode()
	}
	if src.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	srclist := src.Value.(*RedisList)
	if srclist.Len() == 0 {
		return NilBulkReplyEncode()
	}
	dst, err := d.GetObject(resp.Argv[2])
	if err != nil {
		return NilBulkReplyEncode()
	}
	if dst.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	dstlist := src.Value.(*RedisList)

	value := srclist.RPop()
	dstlist.Push([]string{value})

	return StatusReplyEncode(value)
}

func LInsert(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}
	if object.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	l := object.Value.(*RedisList)
	if l.Len() == 0 {
		return IntReplyEncode(0)
	}
	return IntReplyEncode(l.Insert(resp.Argv[2], resp.Argv[3], resp.Argv[4]))
}
func LLen(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}
	if object.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	l := object.Value.(*RedisList)
	return IntReplyEncode(l.Len())
}
func LRem(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}
	if object.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	l := object.Value.(*RedisList)
	if l.Len() == 0 {
		return IntReplyEncode(0)
	}
	count, err := strconv.Atoi(resp.Argv[2])
	if err != nil {
		return ErrReplyEncode(ErrCommand.Error())
	}
	return IntReplyEncode(l.Rem(count, resp.Argv[3]))
}
func LSet(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return IntReplyEncode(0)
	}
	if object.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	l := object.Value.(*RedisList)
	index, err := strconv.Atoi(resp.Argv[2])
	if err != nil {
		return ErrReplyEncode(ErrCommand.Error())
	}
	err = l.Set(index, resp.Argv[3])
	if err == nil {
		return StatusOkReplyEncode()
	} else {
		return ErrReplyEncode(err.Error())
	}
}
func LIndex(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	index, err := strconv.Atoi(resp.Argv[2])
	if err != nil {
		return ErrReplyEncode(ErrCommand.Error())
	}
	l := object.Value.(*RedisList)
	if l.Len() == 0 {
		return NilBulkReplyEncode()
	}
	v, err := l.Index(index)
	if err != nil {
		return NilBulkReplyEncode()
	}
	return BulkReplyEncode(v)
}
func LRange(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return NilBulkReplyEncode()
	}
	if object.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	start, err := strconv.Atoi(resp.Argv[2])
	if err != nil {
		return ErrReplyEncode(ErrCommand.Error())
	}
	end, err := strconv.Atoi(resp.Argv[3])
	if err != nil {
		return ErrReplyEncode(ErrCommand.Error())
	}
	l := object.Value.(*RedisList)
	if l.Len() == 0 {
		return NilBulkReplyEncode()
	}
	ret := l.Range(start, end)
	if len(ret) == 0 {
		return NilBulkReplyEncode()
	}

	return ArrayReplyEncode(ret)
}
func LTrim(d *DB, resp *Resp) []byte {
	object, err := d.GetObject(resp.GetKey())
	if err != nil {
		return StatusOkReplyEncode()
	}
	if object.Type != TypeRedisList {
		return ErrReplyEncode(ErrWrongKeyType.Error())
	}
	start, err := strconv.Atoi(resp.Argv[2])
	if err != nil {
		return ErrReplyEncode(ErrCommand.Error())
	}
	end, err := strconv.Atoi(resp.Argv[3])
	if err != nil {
		return ErrReplyEncode(ErrCommand.Error())
	}
	l := object.Value.(*RedisList)
	if l.Len() == 0 {
		return StatusOkReplyEncode()
	}
	err = l.Trim(start, end)
	if err != nil {
		return ErrReplyEncode(err.Error())
	}

	return StatusOkReplyEncode()
}

func FlushDB(d *DB, resp *Resp) []byte {
	d.Objects = &sync.Map{}
	d.Length = 1
	return StatusOkReplyEncode()
}
func FlushAll(d *DB, resp *Resp) []byte {
	d.Server.FlushAll()
	return StatusOkReplyEncode()
}
func DBSize(d *DB, resp *Resp) []byte {
	return IntReplyEncode(int(d.Length))
}

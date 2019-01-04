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

func (d *DB) Del(keys []string) int {
	ret := 0
	for _,v := range keys {
		_,ok := d.Objects.Load(v)
		if !ok {
			continue
		}
		d.Objects.Delete(v)
		ret++
	}
	return ret
}

func (d *DB) Set(argv []string) error {
	var NXFlag,XXFlag bool
	EXValue := -1
	key,value := argv[0],argv[1]
	argc := len(argv)
	if argc > 2 {
		index := 2
		ctime := int(time.Now().UnixNano() / 1e6)
		for index < argc {
			switch strings.ToUpper(argv[index]) {
			case NX:
				NXFlag = true
				index++
			case XX:
				XXFlag = true
				index++
			case PX:
				if index + 1 >= argc {
					return ErrCommand
				}
				t := argv[index + 1]
				r,err:= strconv.Atoi(t)
				if err != nil{
					return ErrCommand
				}
				EXValue = ctime + r
				index+= 2
			case EX:
				if index + 1 >= argc {
					return ErrCommand
				}
				t := argv[index + 1]
				r,err:= strconv.Atoi(t)
				if err != nil{
					return ErrCommand
				}
				EXValue = ctime + r * 1000
				index+= 2
			default:
				return ErrCommand
			}
		}
	}

	_,ok := d.Objects.Load(key)
	if (ok && NXFlag) || (!ok && XXFlag) {
		return ErrCommand
	}

	o := NewObject()
	o.value = value
	o.Type = RedisString
	o.Encoding = RedisEncodingRaw
	o.Name = key
	o.ExpireAt = EXValue

	d.Objects.Store(key,o)

	return nil
}

func (d *DB) Get(key string) (string,error) {
	v,ok := d.Objects.Load(key)
	if !ok {
		return "",ErrKeyNotFound
	}
	object := v.(*Object)
	if !object.CheckTTL() {
		d.Objects.Delete(key)
		return "",ErrKeyNotFound
	}
	if object.Type != RedisString {
		return "",ErrTypeNotMatch
	}
	value := object.value.(string)
	return value,nil
}

func (d *DB) TTL(key string) (int) {
	v, ok := d.Objects.Load(key)
	if !ok {
		return -2
	}
	object := v.(*Object)
	if !object.CheckTTL() {
		d.Objects.Delete(key)
		return -2
	}

	if object.ExpireAt < 0 {
		return object.ExpireAt
	}else {
		t := object.ExpireAt - int(time.Now().UnixNano() / 1e6)
		return int(math.Ceil(float64(t) / 1000.0))
	}

}

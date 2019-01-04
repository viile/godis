package main

import (
	"strconv"
	"strings"
	"sync"
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
		for index < argc {
			switch strings.ToUpper(argv[index]) {
			case NX:
				NXFlag = true
				index++
			case XX:
				XXFlag = true
				index++
			case PX:
				if index + 1 < argc {
					return ErrCommand
				}
				t := argv[index]
				r,err:= strconv.Atoi(t)
				if err != nil{
					return ErrCommand
				}
				EXValue = r
				index+= 2
			case EX:
				if index + 1 < argc {
					return ErrCommand
				}
				t := argv[index]
				r,err:= strconv.Atoi(t)
				if err != nil{
					return ErrCommand
				}
				EXValue = r * 1000
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
	if object.Type != RedisString {
		return "",ErrTypeNotMatch
	}
	value := object.value.(string)
	return value,nil
}

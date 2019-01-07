package main

import (
	"strconv"
	"sync"
)

type RedisHash struct {
	Objects *sync.Map
	Length int
}

func NewRedisHash() *RedisHash {
	return &RedisHash{
		Objects:&sync.Map{},
	}
}

func (h *RedisHash) Set(key string,value string) int {
	_,ok := h.Objects.Load(key)
	if ok {
		h.Objects.Store(key,value)
		return 0
	} else {
		h.Length++
		h.Objects.Store(key,value)
		return 1
	}
}
func (h *RedisHash) SetNX(key string,value string) int {
	_,ok := h.Objects.Load(key)
	if ok {
		return 0
	}

	h.Length++
	h.Objects.Store(key,value)
	return 1
}

func (h *RedisHash) Get(key string) (string,error) {
	obj,ok := h.Objects.Load(key)
	if !ok {
		return "",ErrKeyNotFound
	}
	return obj.(string),nil
}

func (h *RedisHash) Del(keys []string) int {
	ret :=0
	for _,v := range keys {
		_,ok := h.Objects.Load(v)
		if !ok{
			continue
		}
		h.Objects.Delete(v)
		h.Length--
		ret++
	}
	return ret
}

func (h *RedisHash) Exists(key string) int {
	_,ok := h.Objects.Load(key)
	if !ok {
		return 0
	}
	return 1
}

func (h *RedisHash) GetAll() []string {
	var ret []string
	h.Objects.Range(func(k, v interface{}) bool {
		ret = append(ret,k.(string))
		ret = append(ret,v.(string))
		return true
	})

	return ret
}
func (h *RedisHash) Keys() []string {
	var ret []string
	h.Objects.Range(func(k, v interface{}) bool {
		ret = append(ret,k.(string))
		return true
	})

	return ret
}
func (h *RedisHash) Vals() []string {
	var ret []string
	h.Objects.Range(func(k, v interface{}) bool {
		ret = append(ret,v.(string))
		return true
	})

	return ret
}

func (h *RedisHash) Len() int {
	return h.Length
}

func (h *RedisHash) IncrBy(key string,value int) (int,error) {
	obj,ok := h.Objects.Load(key)
	if ok {
		o := obj.(string)
		oo,err := strconv.Atoi(o)
		if err != nil {
			return 0,ErrWrongKeyType
		}
		ret := oo + value
		h.Objects.Store(key,strconv.Itoa(ret))
		return ret,nil
	} else {
		h.Length++
		h.Objects.Store(key,strconv.Itoa(value))
		return value,nil
	}
}

func (h *RedisHash) IncrByFloat(key string,value float64) (float64,error) {
	obj,ok := h.Objects.Load(key)
	if ok {
		o := obj.(string)
		oo,err := strconv.ParseFloat(o, 64)
		if err != nil {
			return 0,ErrWrongKeyType
		}
		ret := oo + value
		h.Objects.Store(key,strconv.FormatFloat(ret, 'f', -1, 64))
		return ret,nil
	} else {
		h.Length++
		h.Objects.Store(key,strconv.FormatFloat(value, 'f', -1, 64))
		return value,nil
	}
}
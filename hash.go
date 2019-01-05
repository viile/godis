package main

import (
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
	ret := make([]string,h.Length * 2)
	h.Objects.Range(func(k, v interface{}) bool {
		ret = append(ret,k.(string))
		ret = append(ret,v.(string))
		return true
	})

	return ret
}

func (h *RedisHash) Len() int {
	return h.Length
}
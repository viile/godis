package main

import "strconv"

type RedisString struct {
	value string
	length int
}

func NewRedisString(v string) *RedisString{
	return &RedisString{
		value: v,
		length:len(v),
	}
}

func (s *RedisString) Append(str string) int {
	s.value = s.value + str
	s.length = len(s.value)
	return s.length
}

func (s *RedisString) Decr() (int,error) {
	v,err := strconv.Atoi(s.value)
	if err != nil {
		return 0,ErrValueType
	}
	value := v - 1
	s.value = string(value)
	s.length = len(s.value)
	return value,nil
}
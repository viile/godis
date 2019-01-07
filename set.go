package main

import (
	"github.com/deckarep/golang-set"
)

type RedisSet struct {
	Objects mapset.Set
}

func NewRedisSet() *RedisSet {
	return &RedisSet{
		Objects: mapset.NewSet(),
	}
}

func (s *RedisSet) Add(vals []string) int {
	ret := 0
	for _, v := range vals {
		ok := s.Objects.Add(v)
		if ok {
			ret++
		}
	}
	return ret
}

func (s *RedisSet) Card() int {
	return s.Objects.Cardinality()
}

func (s *RedisSet) Diff(o *RedisSet) *RedisSet {
	return &RedisSet{
		Objects: s.Objects.Difference(o.Objects),
	}
}

func (s *RedisSet) Inter(o *RedisSet) *RedisSet {
	return &RedisSet{
		Objects: s.Objects.Intersect(o.Objects),
	}
}
func (s *RedisSet) Union(o *RedisSet) *RedisSet {
	return &RedisSet{
		Objects: s.Objects.Union(o.Objects),
	}
}

func (s *RedisSet) IsMember(value string) bool {
	return s.Objects.Contains(value)
}

func (s *RedisSet) Members() []string {
	ret := make([]string, 0)
	s.Objects.Each(func(i interface{}) bool {
		ret = append(ret, i.(string))
		return false
	})

	return ret
}

func (s *RedisSet) Move(d *RedisSet, value string) bool {
	d.Objects.Add(value)
	s.Objects.Remove(value)
	return true
}

func (s *RedisSet) Pop() string {
	return s.Objects.Pop().(string)
}

func (s *RedisSet) Rem(vals []string) int {
	ret := 0
	for _, v := range vals {
		if s.Objects.Contains(v) {
			s.Objects.Remove(v)
			ret++
		}
	}
	return ret
}

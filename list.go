package main

import (
	"container/list"
)

type RedisList struct {
	Objects *list.List
}

func NewRedisList() *RedisList{
	return &RedisList{
		Objects:list.New(),
	}
}

func (l *RedisList) Push(argv []string) int {
	for _,v := range argv{
		l.Objects.PushFront(v)
	}
	return l.Objects.Len()
}

func (l *RedisList) Insert(op string,pivot string,value string) int {
	if l.Objects.Len() == 0 {
		return 0
	}
	for e := l.Objects.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == pivot {
			if op == "BEFORE" {
				l.Objects.InsertBefore(value,e)
			} else {
				l.Objects.InsertAfter(value,e)
			}

			return l.Objects.Len()
		}
	}

	return -1
}

func (l *RedisList) Len() int {
	return l.Objects.Len()
}

func (l *RedisList) LPop() string {
	e := l.Objects.Front()
	l.Objects.Remove(e)
	return e.Value.(string)
}
func (l *RedisList) RPop() string {
	e := l.Objects.Back()
	l.Objects.Remove(e)
	return e.Value.(string)
}

func (l *RedisList) Rem(count int,value string) int {
	ret := 0
	if count > 0 {
		for e := l.Objects.Front(); e != nil; e = e.Next() {
			if e.Value.(string) == value {
				l.Objects.Remove(e)
				ret++
				if ret == count {
					return count
				}
			}
		}
	} else if count < 0 {
		for e := l.Objects.Back(); e != nil; e = e.Prev() {
			if e.Value.(string) == value {
				l.Objects.Remove(e)
				ret--
				if ret == count {
					return -count
				}
			}
		}
	} else {
		for e := l.Objects.Front(); e != nil; e = e.Next() {
			if e.Value.(string) == value {
				l.Objects.Remove(e)
				ret++
			}
		}
	}
	return ret
}

func (l *RedisList) Set(index int,value string) error {
	ret := 0
	if index >= 0 {
		if index >= l.Objects.Len() {
			return ErrIndexOutOfRange
		}
		for e := l.Objects.Front(); e != nil; e = e.Next() {
			if index == ret {
				e.Value = value
				return nil
			}
			ret++
		}
	} else {
		if -index > l.Objects.Len() {
			return ErrIndexOutOfRange
		}
		for e := l.Objects.Back(); e != nil; e = e.Prev() {
			ret--
			if index == ret {
				e.Value = value
				return nil
			}
		}
	}

	return ErrIndexOutOfRange
}

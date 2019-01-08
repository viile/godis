package main

import (
	"container/list"
	"log"
	"strings"
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
func (l *RedisList) RPush(argv []string) int {
	for _,v := range argv{
		l.Objects.PushBack(v)
	}
	return l.Objects.Len()
}

func (l *RedisList) Insert(op string,pivot string,value string) int {
	if l.Objects.Len() == 0 {
		return 0
	}
	for e := l.Objects.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == pivot {
			if strings.ToUpper(op) == "BEFORE" {
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
	if count >= 0 {
		for e := l.Objects.Front(); e != nil; e = e.Next() {
			log.Println(e.Value)
			if e.Value.(string) == value {
				next := e.Next()
				l.Objects.Remove(e)
				e = next
				ret++
				if ret == count {
					return count
				}
			}
		}
	} else {
		for e := l.Objects.Back(); e != nil; e = e.Prev() {
			if e.Value.(string) == value {
				next := e.Prev()
				l.Objects.Remove(e)
				e = next
				ret--
				if ret == count {
					return -count
				}
			}
		}
	}
	return ret
}

func (l *RedisList) Index(index int) (string,error) {
	ret := 0
	if index >= 0 {
		for e := l.Objects.Front(); e != nil; e = e.Next() {
			if ret == index {
				return e.Value.(string),nil
			}
			ret++
		}
	} else {
		for e := l.Objects.Back(); e != nil; e = e.Prev() {
			ret--
			if ret == index {
				return e.Value.(string),nil
			}
		}
	}

	return "",ErrIndexOutOfRange
}

func (l *RedisList) Range(start int,end int) []string {
	index := -1
	ret := make([]string,0)
	llen := l.Objects.Len()
	/* convert negative indexes */
	if start < 0 {
		start = llen+start
	}
	if end < 0 {
		end = llen+end
	}
	if start < 0 {
		start = 0
	}

	/* Invariant: start >= 0, so this test will be true when end < 0.
	 * The range is empty when start > end or start >= length. */
	if start > end || start >= llen {
		return nil
	}
	if end >= llen {
		end = llen-1
	}

	for e := l.Objects.Front(); e != nil; e = e.Next() {
		index++
		if index < start || index > end {
			continue
		}
		ret = append(ret,e.Value.(string))
	}

	return ret
}

func (l *RedisList) Trim(start int,end int) error {
	index := 0
	llen := l.Objects.Len()
	/* convert negative indexes */
	if start < 0 {
		start = llen+start
	}
	if end < 0 {
		end = llen+end
	}
	if start < 0 {
		start = 0
	}

	/* Invariant: start >= 0, so this test will be true when end < 0.
	 * The range is empty when start > end or start >= length. */
	if start > end || start >= llen {
		l.Objects.Init()
		return nil
	}
	if end >= llen {
		end = llen-1
	}

	for e := l.Objects.Front(); e != nil; e = e.Next() {
		if index < start || index > end {
			next := e.Next()
			l.Objects.Remove(e)
			e = next
		}
		index++
	}

	return nil
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

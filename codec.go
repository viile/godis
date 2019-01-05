package main

import (
	"strconv"
)

type Codec struct {
	Buf           []byte
	Argv          []string
	Argc          int
	LastBulkIndex int
	status        int
}

type Resp struct {
	Argc int
	Argv []string
}

func (r *Resp) GetCommand() string {
	return r.Argv[0]
}
func (r *Resp) GetByIndex(start int,end int) []string {
	return r.Argv[start:end]
}
func (r *Resp) GetKey() string {
	return r.Argv[1]
}

func NewCodec() *Codec {
	return &Codec{
		Buf:    make([]byte, 0),
		Argv:    make([]string, 0),
	}
}


func (c *Codec) init() {
	c.Buf = make([]byte, 0)
	c.Argv = make([]string, 0)
	c.Argc = 0
	c.LastBulkIndex = 0
	c.status = InitParser
}

// Decode .
func (c *Codec) Decode(buf *[]byte) (Resps []*Resp) {
	Resps = make([]*Resp,0)
	Buf := append(c.Buf,*buf...)
	Length := len(Buf)
	var BulkLength = 0
	var index = 0
	for index < Length {
		switch c.status {
		case InitParser:
			switch RespReply(Buf[index]) {
			case ArrayReply:
				c.status = ArrayParser
				index++
			default:
				c.init()
				return
			}
		case ArrayParser:
			if Buf[index] != '\r' {
				c.Argc = (c.Argc * 10) + int(Buf[index]) - 48
				index++
			} else {
				c.status = BulkLengthParser
				index = index + 2
				c.LastBulkIndex = index
			}
		case BulkLengthParser:
			// bulk check
			if RespReply(Buf[index]) == BulkLengthReply {
				index++
				c.status = BulkStringParser
			} else {
				c.init()
				return
			}
		case BulkStringParser:
			if Buf[index] != '\r' {
				BulkLength = (BulkLength * 10) + int(Buf[index]) - 48
				index++
			} else {
				if BulkLength < 0 {
					c.init()
					return
				} else {
					BulkStartIndex := index + 2
					BulkEndIndex := BulkStartIndex + BulkLength
					if BulkEndIndex > Length {
						// 长度不够 等待下一个包
						c.Buf = Buf[c.LastBulkIndex:]
						return
					}

					c.Argv = append(c.Argv, string(Buf[BulkStartIndex:BulkEndIndex]))
					index = BulkEndIndex + 2
					c.LastBulkIndex = index
					BulkLength = 0
					c.status = BulkLengthParser

					if c.Argc == len(c.Argv) {
						Resps = append(Resps,&Resp{Argc:c.Argc,Argv:c.Argv})
						c.init()
					}
				}
			}
		default:
			c.init()
			return
		}
	}
	return
}


func ErrReplyEncode(str string) []byte {
	ret := make([]byte,0)
	ret = append(ret,byte(ErrorReply))
	ret = append(ret,[]byte(str)...)
	ret = append(ret,[]byte{13,10}...)
	return ret
}

func IntReplyEncode(i int) []byte {
	ret := make([]byte,0)
	ret = append(ret,byte(IntReply))
	ret = append(ret,[]byte(Int2bytes(i))...)
	ret = append(ret,[]byte{13,10}...)
	return ret
}

func StatusReplyEncode(str string) []byte {
	ret := make([]byte,0)
	ret = append(ret,byte(StatusReply))
	ret = append(ret,[]byte(str)...)
	ret = append(ret,[]byte{13,10}...)
	return ret
}

func BulkReplyEncode(str string) []byte {
	ret := make([]byte,0)
	ret = append(ret,byte(BulkLengthReply))
	ret = append(ret,Int2bytes(len(str))...)
	ret = append(ret,[]byte{13,10}...)
	ret = append(ret,[]byte(str)...)
	ret = append(ret,[]byte{13,10}...)
	return ret
}

func NilBulkReplyEncode() []byte {
	return []byte{36,45,49,13,10}
}
func StatusOkReplyEncode() []byte {
	return []byte{43,79,75,13,10}
}

func ArrayReplyEncode(strs []string) []byte {
	ret := make([]byte,0)
	ret = append(ret,byte(ArrayReply))
	ret = append(ret,Int2bytes(len(strs))...)
	ret = append(ret,[]byte{13,10}...)
	for _,v := range strs {
		ret = append(ret,BulkReplyEncode(v)...)
	}
	return ret
}

func Int2bytes(i int) []byte {
	s := strconv.Itoa(i)
	r := []byte(s)
	return r
}
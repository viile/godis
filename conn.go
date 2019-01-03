package main

import (
	"bytes"
	"context"
	"net"
)

// Conn wrap net.Conn
type Conn struct {
	rawConn   net.Conn
	sendCh    chan []byte
	done      chan error
	name      string
	messageCh chan *[]byte
}

// GetName Get conn name
func (c *Conn) GetName() string {
	return c.name
}

// NewConn create new conn
func NewConn(c net.Conn) *Conn {
	conn := &Conn{
		rawConn:   c,
		sendCh:    make(chan []byte, 100),
		done:      make(chan error),
		messageCh: make(chan *[]byte, 100),
	}

	conn.name = c.RemoteAddr().String()

	return conn
}

// Close close connection
func (c *Conn) Close() {
	c.rawConn.Close()
}

// SendMessage send message
func (c *Conn) SendMessage(msg *Message) error {
	pkg := []byte{43, 79, 75, 13, 10}
	c.sendCh <- pkg
	return nil
}

// writeCoroutine write coroutine
func (c *Conn) writeCoroutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case pkt := <-c.sendCh:
			if pkt == nil {
				continue
			}

			if _, err := c.rawConn.Write(pkt); err != nil {
				c.done <- err
			}
		}
	}
}

// readCoroutine read coroutine
func (c *Conn) readCoroutine(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 读取数据
			buf := make([]byte,MaxReadSize)
			n,err := c.rawConn.Read(buf)
			if err != nil {
				c.done <- err
				continue
			}
			if n == 0 {
				continue
			}
			r := bytes.TrimRight(buf, "\x00")
			c.messageCh <- &r
		}
	}
}

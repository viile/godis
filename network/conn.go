package network

import (
	"context"
	"io"
	"net"
)

// Conn wrap net.Conn
type Conn struct {
	sid       string
	rawConn   net.Conn
	sendCh    chan []byte
	done      chan error
	name      string
	messageCh chan *Message
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
		messageCh: make(chan *Message, 100),
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
	pkg, err := Encode(msg)
	if err != nil {
		return err
	}
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
			databuf := make([]byte, 8)
			l, err := io.ReadFull(c.rawConn, databuf)
			if err != nil {
				c.done <- err
				continue
			}
			if l == 0 {
				continue
			}

			// 解码
			msg, err := Decode(databuf)
			if err != nil {
				c.done <- err
				continue
			}

			c.messageCh <- msg
		}
	}
}

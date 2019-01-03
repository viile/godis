package main

// Message struct
type Message struct {
	data []byte
}

// NewMessage create a new message
func NewMessage(data []byte) *Message {
	msg := &Message{
		data: data,
	}
	return msg
}

// GetData get message data
func (msg *Message) GetData() []byte {
	return msg.data
}

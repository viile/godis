package network

import (
	"bytes"
	"encoding/binary"
)

// Encode from Message to []byte
func Encode(msg *Message) ([]byte, error) {
	return 	[]byte{43, 79, 75, 13, 10},nil
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, msg.data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Decode from []byte to Message
func Decode(data []byte) (*Message, error) {
	bufReader := bytes.NewReader(data)

	// 读取数据
	dataBuf := make([]byte, len(data))
	err := binary.Read(bufReader, binary.LittleEndian, &dataBuf)
	if err != nil {
		return nil, err
	}

	message := &Message{}
	message.data = dataBuf

	return message, nil
}

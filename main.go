package main

import (
	"fmt"
)

func main() {
	host := "127.0.0.1:18787"

	ss, err := NewSocketService(host)
	if err != nil {
		return
	}

	ss.RegMessageHandler(HandleMessage)
	ss.RegConnectHandler(HandleConnect)
	ss.RegDisconnectHandler(HandleDisconnect)

	gs := NewGodisServer()
	_ = gs
	ss.Serv()
}

func HandleMessage(s *Session, msg *Message) {
	//fmt.Println("receive msgID:", msg)
	fmt.Println(msg.GetData())
	s.GetConn().SendMessage(nil)
}

func HandleDisconnect(s *Session, err error) {
	fmt.Println(s.GetConn().GetName() + " lost.")
}

func HandleConnect(s *Session) {
	fmt.Println(s.GetConn().GetName() + " connected.")
}

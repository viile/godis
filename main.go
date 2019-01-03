package main

import (
	"fmt"
)

func main() {
	host := "127.0.0.1:18787"

	ss, err := NewServer(host)
	if err != nil {
		return
	}

	ss.RegMessageHandler(HandleMessage)
	ss.RegConnectHandler(HandleConnect)
	ss.RegDisconnectHandler(HandleDisconnect)

	ss.Run()
}
// HandleMessage .
func HandleMessage(s *Session, buf *[]byte) {
	//fmt.Println("receive msgID:", msg)
	fmt.Println(buf)
	s.GetConn().SendMessage(nil)
}
// HandleDisconnect .
func HandleDisconnect(s *Session, err error) {
	fmt.Println(s.GetConn().GetName() + " lost.")
}
// HandleConnect .
func HandleConnect(s *Session) {
	fmt.Println(s.GetConn().GetName() + " connected.")
	s.GetConn().SendMessage(nil)
}

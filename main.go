package main

import (
	"flag"
	"log"
)

var Host string

func main() {
	flag.StringVar(&Host, "h", "0.0.0.0:18787", "server host")
	flag.Parse()

	s, err := NewServer(Host)
	if err != nil {
		log.Println(err)
		return
	}

	s.Run()
}


package main

func main() {
	host := "127.0.0.1:18787"

	s, err := NewServer(host)
	if err != nil {
		return
	}

	s.Run()
}


package main

func main() {
	host := "0.0.0.0:18787"

	s, err := NewServer(host)
	if err != nil {
		return
	}

	s.Run()
}


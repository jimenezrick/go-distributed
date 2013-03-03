package main

import "time"

var conn Conn

func main() {
	Start("localhost:8000")
	conn = Connect("localhost:8001")
}

func loop() {
	for {
		conn.out <- []byte("ping")
		time.Sleep(time.Second)
	}
}

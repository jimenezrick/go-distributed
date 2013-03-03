package main

import "net"







var connections = make(map[string] Conn)


type Conn struct {
	tag string
	sock net.Conn
	in <-chan []byte
	out chan<- []byte
}




func sender(Conn conn, out <-chan []byte) {
	for {
		msg := <-out
		_, err := conn.sock.Write(msg)
		if err != nil {
		}
	}
}

func receiver(net.Conn conn, in chan<- []byte) {
}




func Connect() (chan<- []byte, <-chan []byte) {




	in := make(chan []byte, 10)
	out := make(chan []byte, 10)

	go receiver(conn, in)
	go sender(conn, out)

	return in, out
}




func connect(host string) {
	conn, err := net.Dial("tcp", host + ":8000")
	if err != nil {
		panic(err)
	}

	c, ok := connections[host]
	if ok {
		conn.Close()
	}
}

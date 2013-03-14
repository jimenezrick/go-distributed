package main

import (
	"net"
	"sync"
	"time"
)

type mesh struct {
	net.Listener
	connsMtx sync.Mutex
	connsMap map[string] conn // XXX: Use Skip List with tag + node as key
}

type conn struct {
	net.Conn
	in, out chan msg
}

type msg []byte

func newConn(sock net.Conn) conn {
	return conn{sock, make(chan msg), make(chan msg)}
}

func New(laddr string) (*mesh, error) {
	lis, err := net.Listen("tcp", laddr)
	if err != nil {
		return nil, err
	}

	m := &mesh{lis, sync.Mutex{}, make(map[string] conn)}
	go m.acceptor()
	return m, nil
}

func (m *mesh) acceptor() {
	defer m.Close()

	for {
		conn, err := m.Accept()
		if err != nil {
			panic(err)
		}
		go m.handleConn(conn)
	}
}

func (m *mesh) initConn(conn net.Conn) {
	//
	// XXX: Read connection tag, do not include in the messages
	//
	tag := make([]byte, 128)
	conn.Read(tag)
	m.registerConn(string(tag), conn)
}

func (m *mesh) handleConn(conn net.Conn) {
	m.initConn(conn)
	//
	// XXX
	//
}

func (m *mesh) registerConn(tag string, sock net.Conn) {
	defer m.connsMtx.Unlock()
	m.connsMtx.Lock()

	// XXX: Insert pointer to connections in the map so we can do
	// this in-place modification of the socket?
	if conn, ok := m.connsMap[tag]; ok {
		if conn.Conn != nil {
			conn.Conn = sock
		} else {
			panic(conn.Conn)
		}
	} else {
		m.connsMap[tag] = newConn(sock)
	}
}

func (m *mesh) Connect(tag, addr string) {
	go m.connect(addr, addr)
}

func (m *mesh) connect(tag, addr string) {
	var (
		sock net.Conn
		err error
	)

	for {
		// TODO: Log, parametrize sleep
		sock, err = net.Dial("tcp", addr)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	m.registerConn(tag, sock)
}

func (m *mesh) OutQueue(tag, addr string) <-chan msg {
	defer m.connsMtx.Unlock()
	m.connsMtx.Lock()

	conn, ok := m.connsMap[tag]
	if !ok {
		conn = newConn(nil)
		m.connsMap[tag] = conn
		m.Connect(tag, addr)
	}
	return conn.out
}

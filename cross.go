package main

import (
	"net"
	"sync"
	"time"
)

type tag string

type msg []byte

type mesh struct {
	net.Listener
	connsMtx sync.Mutex
	connsMap map[tag] *conn // XXX: Use Skip List
}

type conn struct {
	net.Conn
	in, out chan msg
}

func newConn(sock net.Conn) *conn {
	return &conn{sock, make(chan msg), make(chan msg)}
}

func New(laddr string) (*mesh, error) {
	lis, err := net.Listen("tcp", laddr)
	if err != nil {
		return nil, err
	}

	m := &mesh{lis, sync.Mutex{}, make(map[tag] *conn)}
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
		go m.connHandler(conn)
	}
}

func (m *mesh) initConn(conn net.Conn) {
	//
	// XXX: Read connection tag, do not include in the messages
	//
	t := make([]byte, 128)
	conn.Read(t)
	m.registerConn(tag(t), conn)
}

func (m *mesh) connHandler(conn net.Conn) {
	m.initConn(conn)
	//
	// XXX
	//
}

func (m *mesh) registerConn(tag tag, sock net.Conn) {
	defer m.connsMtx.Unlock()
	m.connsMtx.Lock()

	if conn, ok := m.connsMap[tag]; ok {
		if conn.Conn == nil {
			conn.Conn = sock
		} else {
			sock.Close()
		}
	} else {
		m.connsMap[tag] = newConn(sock)
	}
	//
	// XXX: Spawn reader/writer goroutines
	//
}

func (m *mesh) Connect(t tag, addr string) {
	go m.connecter(t, addr)
}

func (m *mesh) connecter(t tag, addr string) {
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

	m.registerConn(t, sock)
}

func (m *mesh) OutQueue(t tag, addr string) <-chan msg {
	defer m.connsMtx.Unlock()
	m.connsMtx.Lock()

	conn, ok := m.connsMap[t]
	if !ok {
		conn = newConn(nil)
		m.connsMap[t] = conn
		m.Connect(t, addr)
	}
	return conn.out
}

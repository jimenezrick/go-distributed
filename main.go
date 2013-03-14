package main

func main() {
	m, err := New("localhost:8000");
	if err != nil {
		panic(err)
	}

	c := m.OutQueue("foo", "localhost:8001")
	msg := <-c
	println(msg)
}

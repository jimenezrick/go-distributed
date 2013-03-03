import "mesh"



int Main() {

	m := mesh.New()
	m:Start()


	m.Connect("a.com")

	// /stats
	// /stats/net
	// /stats/cpu
	// /stats/mem

	c := m.Subscribe("/info")
	c1 := m.Subscribe("/stats")
	c2 := m.SubscribeNode("foo", "/stats")

	m.Send(node, "/jfhdskjfh", data)
	select {
	case msg := <-c:
		....
	case <-time.After(1000):
		...
	}



	m.Send(node, "/alerts", data)
	/*
	m.SendAll("/alerts", data)
	*/


	p ! Msg
	receive
		{node, Res} ->
			<<correct path>>
			.....
			.....
	after 1000 ->
		<<timeout error>>
		.....
		.....
}

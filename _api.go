package "mesh"

/*
   Me gusta el nombre de "mesh" para el paquete, porque es corto y representa un
   poco la idea de lo que haria la libreria. Pero vamos, que es un tema totalmente
   abierto y no es el definitivo.
 */

/*
   Un tema que no tengo claro del todo es lo de usar un solo socket. Es lo mas sencillo
   y elegante, pero si un canal se llena, puede bloquear todo la comunicacion entre 2
   nodos y es una cosa que en Erlang no esta resuelta y estaria bien hacer algo al respecto.
 */

int Main() {
	// Arrancamos la libreria, principalmente la gorutina aceptora
	// e iniciaizamos las estructuras de datos.
	m, err := mesh.New()
	m.Start()

	/*
	   La libreria deberia dejar parametrizar su comportamiento:
	     - Tamaños de colas
	     - Comportamiento cuando los canales estan llenos
	     - Tiempo de los timeouts
	 */

	// Pedimos un intento de conexion asincrono.
	m.Connect("a.com")
	// Pedimos un intento de conexion sincrono. Esperamos hasta conseguir
	// una conexion satisfactoria con el nodo dado. Pero eso no impide que
	// tras conseguir la conexion, la comunicacion se rompa. Pero de todos
	// modos a veces puede ser comodo esperar a que la conexion se establezca
	// para no perder muchos mensajes.
	err := m.WaitConnected("b.com")

	/*
	   La etiqueta asociada a un mensaje es una simple ristra de bytes. Se
	   Puede interpretar como un string. La subscripcion se basa en si el
	   prefijo pedido encaja con la etiquita del paquete. El formato
	   de la etiqueta es totalmente libre, pero puede ser buena ida (no obligatoria)
	   usar un espacio de nombres similar al espacio de nombres del FS de UNIX:

	     /stats
	     /stats/net
	     /stats/cpu
	     /stats/mem
	 */

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

	// El envio encola el mensaje en un canal interno que la gorutina enviadora
	// leera de el. Se puede hacer una operacion no bloqueante de envio por si
	// dicho canal esta lleno. O igualmente parametrizar que sucede si se intenta
	// encolar en un canal que esta lleno.
	m.Send(node, "/alerts", data)
	m.SendAll("/alerts", data)
	m.SendNonBlocking("/alerts", data)

	/*
	   Otra cuestion acerca del envio, es si en vez de tener un metodo Send() decidimos
	   tambien exponer canales. Del mismo modo que tenemos canales para recibir. A lo mejor
	   es mas elegante darle al usuario canales para enviar, que seria lo simetrico a tener canales
	   para recibir. Por tanto tener un metodo que crea el canal y la gorutina enviadora asociada.
	*/
	// Crea un canal por el que mandar cosas asociadas a esa etiqueta y la gorutina asociada
	// enviadora que estara pendiente de enviar esos datos por el socket.
	c := m.Publish("/alerts")
	c <- msg

	/*
		p ! Msg
		receive
			{node, Res} ->
				<<correct path>>
				.....
		after 1000 ->
			<<timeout error>>
			.....
	 */
}

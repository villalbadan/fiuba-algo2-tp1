package lista

type Lista[T any] interface {

	//Devuelve true si hay elementos en la lista, devuelve false en caso contrario.
	EstaVacia() bool

	//Agrega un elemento al primer lugar de la lista.
	InsertarPrimero(T)

	//Agrega un elemento al último lugar de la lista.
	InsertarUltimo(T)

	//Elimina el primer elemento que se encuentra en la lista,
	//si la misma está vacía entra en pánico con el mensaje "La lista esta vacia".
	BorrarPrimero() T

	//Devuelve el primer elemento de la lista.
	//Si la lista está vacía entra en pánico con el mensaje "La lista esta vacia".
	VerPrimero() T

	//Devuelve el último elemento de la lista.
	//Si la lista está vacía entra en pánico con el mensaje "La lista esta vacia"
	VerUltimo() T

	//Devuelve el largo de la lista.
	Largo() int

	//Itera por todos los elementos de la lista mientras que la función recibida sea verdadera,
	//En caso contrario termina la iteración.
	Iterar(func(T) bool)

	// Iterador apunta al primer elemento de la lista.
	Iterador() IteradorLista[T]
}

type IteradorLista[T any] interface {

	//VerActual devuelve el elemento donde está parado el iterador actualmente.
	//Si ya ha iterado todos los elementos, entra en pánico con ("El iterador termino de iterar")
	VerActual() T

	//HaySiguiente devuelve verdadero si el nodo al que se encuentra apuntando el interador no esta vacio.
	//Devuelve falso en caso contrario (fin de la lista).
	HaySiguiente() bool

	//Siguiente devuelve el elemento donde estaba parado el iterador y avanza a la siguiente posición de la lista.
	//Si ya ha iterado todos los elementos, entra en pánico con ("El iterador termino de iterar")
	Siguiente() T

	//Inserta el elemento recibido en la posición anterior a la que se encuentra el iterador.
	//El mismo pasa a apuntar al elemento insertado.
	Insertar(T)

	//Borrar elimina el elemento al cual estaba apuntando el iterador
	//y pasa a apuntar al elemento siguiente del eliminado.
	//Si ya ha iterado todos los elementos, entra en pánico con ("El iterador termino de iterar")
	Borrar() T
}

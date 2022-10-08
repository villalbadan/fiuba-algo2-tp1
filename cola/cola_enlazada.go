package cola

type colaEnlazada[T any] struct {
	primero *nodoCola[T]
	ultimo  *nodoCola[T]
}

type nodoCola[T any] struct {
	dato T
	prox *nodoCola[T]
}

func CrearColaEnlazada[T any]() Cola[T] {
	return new(colaEnlazada[T])
}

func crearNuevoNodo[T any](data T) *nodoCola[T] {
	return &nodoCola[T]{dato: data}
}

// EstaVacia devuelve verdadero si la cola no tiene elementos encolados, false en caso contrario.
func (cola colaEnlazada[T]) EstaVacia() bool {
	return cola.primero == nil
}

// VerPrimero obtiene el valor del primero de la cola. Si está vacía, entra en pánico con un mensaje
// "La cola esta vacia".
func (cola colaEnlazada[T]) VerPrimero() T {
	if cola.EstaVacia() {
		panic("La cola esta vacia")
	}
	return cola.primero.dato

}

// Encolar agrega un nuevo elemento a la cola, al final de la misma.
func (cola *colaEnlazada[T]) Encolar(data T) {
	nuevoNodo := crearNuevoNodo[T](data)

	if cola.EstaVacia() {
		cola.primero = nuevoNodo
	} else {
		cola.ultimo.prox = nuevoNodo
	}

	cola.ultimo = nuevoNodo
}

// Desencolar saca el primer elemento de la cola. Si la cola tiene elementos, se quita el primero de la misma,
// y se devuelve ese valor. Si está vacía, entra en pánico con un mensaje "La cola esta vacia".
func (cola *colaEnlazada[T]) Desencolar() T {
	if cola.EstaVacia() {
		panic("La cola esta vacia")
	}

	dataDesencolada := cola.primero.dato
	if cola.primero.prox == nil {
		cola.ultimo = nil
	}
	cola.primero = cola.primero.prox

	return dataDesencolada
}

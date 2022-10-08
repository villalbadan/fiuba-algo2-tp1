package lista

type listaEnlazada[T any] struct {
	primero *nodoLista[T]
	ultimo  *nodoLista[T]
	largo   int
}

type nodoLista[T any] struct {
	dato T
	prox *nodoLista[T]
}

type iteradorExterno[T any] struct {
	lista    *listaEnlazada[T]
	actual   *nodoLista[T]
	anterior *nodoLista[T]
}

func CrearListaEnlazada[T any]() Lista[T] {
	return new(listaEnlazada[T])
}

func crearNuevoNodo[T any](elemento T) *nodoLista[T] {
	return &nodoLista[T]{dato: elemento}
}

// PRIMITIVAS DE LISTA ENLAZADA ----------------------------------------------------

func (lista *listaEnlazada[T]) EstaVacia() bool {
	return lista.largo == 0
}

func (lista *listaEnlazada[T]) InsertarPrimero(elemento T) {
	nuevoNodo := crearNuevoNodo(elemento)

	if lista.EstaVacia() {
		lista.ultimo = nuevoNodo
	} else {
		nuevoNodo.prox = lista.primero
	}

	lista.primero = nuevoNodo
	lista.largo++
}

func (lista *listaEnlazada[T]) InsertarUltimo(elemento T) {
	nuevoNodo := crearNuevoNodo(elemento)

	if lista.EstaVacia() {
		lista.primero = nuevoNodo
	} else {
		lista.ultimo.prox = nuevoNodo
	}

	lista.ultimo = nuevoNodo
	lista.largo++
}

func (lista *listaEnlazada[T]) BorrarPrimero() T {
	if lista.EstaVacia() {
		panic("La lista esta vacia")
	}

	elementoEliminado := lista.primero.dato
	lista.primero = lista.primero.prox
	lista.largo--
	return elementoEliminado
}

func (lista listaEnlazada[T]) VerPrimero() T {
	if lista.EstaVacia() {
		panic("La lista esta vacia")
	}
	return lista.primero.dato
}

func (lista listaEnlazada[T]) VerUltimo() T {
	if lista.EstaVacia() {
		panic("La lista esta vacia")
	}
	return lista.ultimo.dato
}

func (lista listaEnlazada[T]) Largo() int {
	return lista.largo
}

func (lista listaEnlazada[T]) Iterar(visitar func(T) bool) {
	actual := lista.primero
	for i := 0; i < lista.largo; i++ {
		if !visitar(actual.dato) {
			break
		}
		actual = actual.prox
	}
}

func (lista *listaEnlazada[T]) Iterador() IteradorLista[T] {
	return &iteradorExterno[T]{lista: lista, actual: lista.primero}
}

// PRIMITIVAS DE ITERADOR EXTERNO -------------------------------------------------

func (iter iteradorExterno[T]) VerActual() T {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	return iter.actual.dato
}

func (iter iteradorExterno[T]) HaySiguiente() bool {
	return iter.actual != nil
}

func (iter *iteradorExterno[T]) Siguiente() T {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	iter.anterior = iter.actual
	iter.actual = iter.actual.prox
	return iter.anterior.dato
}

func (iter *iteradorExterno[T]) Borrar() T {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	dato := iter.actual.dato
	if iter.anterior != nil {
		iter.anterior.prox = iter.actual.prox
	} else {
		iter.lista.primero = iter.actual.prox
	}

	if iter.actual.prox == nil {
		iter.lista.ultimo = iter.anterior
	}

	iter.actual = iter.actual.prox
	iter.lista.largo--
	return dato
}

func (iter *iteradorExterno[T]) Insertar(elemento T) {
	nuevoNodo := crearNuevoNodo(elemento)
	nuevoNodo.prox = iter.actual

	if iter.anterior != nil {
		iter.anterior.prox = nuevoNodo
	} else {
		iter.lista.primero = nuevoNodo
	}

	if iter.actual == nil {
		iter.lista.ultimo = nuevoNodo
	}

	iter.actual = nuevoNodo
	iter.lista.largo++
}

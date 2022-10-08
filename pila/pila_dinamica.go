package pila

const (
	_TAMANIO_INICIAL    = 5
	_FACTOR_REDIMENSION = 2
	_AGRANDAR           = 4
)

/* Definición del struct pila proporcionado por la cátedra. */
type pilaDinamica[T any] struct {
	datos    []T
	cantidad int
}

func CrearPilaDinamica[T any]() Pila[T] {
	pila := new(pilaDinamica[T])
	pila.datos = make([]T, _TAMANIO_INICIAL)
	return pila
}

// EstaVacia devuelve verdadero si la pila no tiene elementos apilados, false en caso contrario.
func (pila *pilaDinamica[T]) EstaVacia() bool {
	return pila.cantidad == 0
}

// VerTope obtiene el valor del tope de la pila. Si la pila tiene elementos se devuelve el valor del tope.
// Si está vacía, entra en pánico con un mensaje "La pila esta vacia".
func (pila *pilaDinamica[T]) VerTope() T {
	if pila.EstaVacia() {
		panic("La pila esta vacia")
	}
	return pila.datos[pila.cantidad-1]
}

// Apilar agrega un nuevo elemento a la pila.
func (pila *pilaDinamica[T]) Apilar(elem T) {
	if len(pila.datos) == pila.cantidad {
		pila.redimensionar(pila.cantidad * _FACTOR_REDIMENSION)
	}

	pila.datos[pila.cantidad] = elem
	pila.cantidad++
}

// Desapilar saca el elemento tope de la pila. Si la pila tiene elementos, se quita el tope de la pila, y
// se devuelve ese valor. Si está vacía, entra en pánico con un mensaje "La pila esta vacia".
func (pila *pilaDinamica[T]) Desapilar() T {
	if pila.EstaVacia() {
		panic("La pila esta vacia")
	}

	pila.cantidad--

	if (len(pila.datos) >= pila.cantidad*_AGRANDAR) && (len(pila.datos)/_FACTOR_REDIMENSION >= _TAMANIO_INICIAL) {
		pila.redimensionar(len(pila.datos) / _FACTOR_REDIMENSION)
	}
	return pila.datos[pila.cantidad]
}

func (pila *pilaDinamica[T]) redimensionar(newSize int) {
	temp := pila.datos
	pila.datos = make([]T, newSize)
	copy(pila.datos, temp)
}

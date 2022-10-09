package tp1

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	TDACola "tp1/cola"
	"tp1/errores"
	TDALista "tp1/lista"
	"tp1/votos"
)

const (
	MIN_DNI = 1000000
	MAX_DNI = 100000000
	INIT_PADRON = 100
	INIT_PARTIDOS = 10
)

//  ############### DESHACER ------------------------------------------------------------------------------------------
func deshacerVoto(fila TDACola.Cola[votos.Votante]) {
	errDeshacer := fila.VerPrimero().Deshacer()
	if errDeshacer != nil {
		fmt.Fprintf(os.Stdout, "%s", errDeshacer)
	} else {
		fmt.Fprintf(os.Stdout,"OK")
	}
}

// ############### INGRESAR DNI ---------------------------------------------------------------------------------------
func buscarEnPadron(padron []votos.Votante, dni int) (votos.Votante, error) {
	// busqueda binaria
	// si no esta >> return nil, errores.DNIFueraPadron{}
	return votanteIngresado, nil
}

func controlarDNI(padron []votos.Votante, dni int) (votos.Votante, error) {
	// controlar si no es un numero muy pequeño asi no se busca de más
	if dni <= MIN_DNI || dni >= MAX_DNI {
		return nil, errores.DNIError{}
	}
	return buscarEnPadron(padron, dni)
}

func ingresarDNI(fila TDACola.Cola[votos.Votante], padron []votos.Votante, dni int) {
	votanteIngresado, errIngresando := controlarDNI(padron, dni)
	if errIngresando == nil {
		fila.Encolar(votanteIngresado)
		fmt.Fprintf(os.Stdout,"OK")
	} else {
		fmt.Fprintf(os.Stdout, "%s", errIngresando)
	}
}

// ############### VOTAR ----------------------------------------------------------------------------------------------
func controlarVoto(fila TDACola.Cola[votos.Votante], datos string) struct{}, error { // no va a ser un struct creo
	// tipo TipoVoto
	// alternativa int
	return {tipo: , alternativa: }, nil
}

func votar(fila TDACola.Cola[votos.Votante], datos string) {
	if fila.EstaVacia() {
		fmt.Fprintf(os.Stdout, "%s", errores.FilaVacia{})
	} else {
		// ni idea todavia si lo que devuelve es un string pero vamos a asumir por ahora, puede ser un array de 2d?
		// tambien puede ya estar separado, segun como lo implementemos al leer las lineas
		voto, errVotar := controlarVoto(fila, datos)
		if errVotar == nil {
			fila.VerPrimero().Votar(voto.tipo, voto.alternativa)
		} else {
			fmt.Fprintf(os.Stdout, "%s", errVotar)
		}
	}
}

// ############### FIN-VOTO  ------------------------------------------------------------------------------------------
func sumarVoto(voto votos.Voto, partidos []votos.Partido, candidaturas []votos.TipoVoto)  {
	for i, _ := range candidaturas {
		partidos[voto.VotoPorTipo[i]].VotadoPara(candidaturas[i])
	}
}

func finalizarVoto(fila TDACola.Cola[votos.Votante], partidos []votos.Partido, impugnados TDALista.Lista[votos.Voto], candidaturas []votos.TipoVoto) {
	voto, errFinalizar := fila.VerPrimero().FinVoto()
	if errFinalizar != nil {
		fmt.Fprintf(os.Stdout, "%s", errFinalizar)
		impugnados.InsertarUltimo(voto)
	} else {
		sumarVoto(voto, partidos, candidaturas)
	}
	fila.Desencolar()
}

// ############### ----------------------------------------------------------------------------------------------------

func prepararLista(lista []votos.Partido, archivoLista string) {
	lista[0] = votos.CrearVotosEnBlanco()

	archivo, err := os.Open(archivoLista)
	defer archivo.Close()

	s := bufio.NewScanner(archivo)
	for s.Scan() {
		dividirLinea := strings.Split(s.Text(), ",")
		partido := votos.CrearPartido(dividirLinea[0], dividirLinea[1:])
		lista = append(lista, partido)
	}

	err = s.Err()
	if err != nil {
		fmt.Println(err)
	}

}

func leerPadron(archivoPadron string) []int {

	temp := make([]int, INIT_PADRON)
	archivo, err := os.Open(archivoPadron)
	defer archivo.Close()

	s := bufio.NewScanner(archivo)
	for s.Scan() {
		linea, _ := strconv.Atoi(s.Text())
		temp = append(temp, linea)
	}
	err = s.Err()
	if err != nil {
		fmt.Println(err)
	}

	return temp
}

func prepararPadron(padron []votos.Votante, archivoPadron string) {
	// Ordenar padron en un array para despues hacer busqueda binaria (en el caso del padron)
	// Ver si podemos no leer el padron 2 veces (una hace el array simple y la otra lo hace con el struct entero pero ya ordenado)
	//y ordenar directamente el struct
	// vi que hay un par de maneras con sort.Slice pero no puedo probar que simplemente funcione asi que no me meti ahi
	temp := leerPadron(archivoPadron)
	sort.Ints(temp)
	for i := range temp {
		padron = append(padron, votos.CrearVotante(temp[i]))
	}
}

func prepararMesa(archivoLista, archivoPadron string) ([]votos.Partido, []votos.Votante) {
	// estructuras que vamos a usar, puse los valores de las const como placeholder pero habria que ver cuantos partidos/dni
	//trae cada archivo de prueba y ahi hacer el array? porque en caso de un archivo de 300mil va a redimensionar banda
	// no se que conviene, sobretodo en el padron, la lista de partidos suele ser corta
	padron := make([]votos.Votante, INIT_PADRON)
	lista := make([]votos.Partido, INIT_PARTIDOS)
	// leer archivos
	prepararPadron(padron, fmt.Sprintf("%s.txt", archivoPadron))
	prepararLista(lista, fmt.Sprintf("%s.csv",archivoLista))
	return lista, padron
}

func inicializar(args []string) bool {
	// parametros correctos
	if len(args) < 2 {
		fmt.Fprintf(os.Stdout, "%s", errores.ErrorParametros{})
		return false
	}

	// archivos existen
	_, err1 := os.Stat(args[0])
	_, err2 := os.Stat(args[1])
	if err2 != nil || err1 != nil {
		fmt.Fprintf(os.Stdout, "%s", errores.ErrorLeerArchivo{})
		return false
	}
	return true
}

func main() {
	var (
		padron     []votos.Votante
		partidos   []votos.Partido
		candidaturas = []votos.TipoVoto{votos.PRESIDENTE, votos.GOBERNADOR, votos.GOBERNADOR}
		impugnados = TDALista.CrearListaEnlazada[votos.Voto]()
		// iba a hacer un array para impugnados y dar como resultado el len del array
		// pero siento que iterar la lista al final va a ser menos costoso que redimensionar tantas veces?
	)

	if inicializar(os.Args[1:]) {
		partidos, padron = prepararMesa(os.Args[1], os.Args[2])
		// cola de votantes
		fila := TDACola.CrearColaEnlazada[votos.Votante]()

		// lectura stdin
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			// leer linea
			// separar en args
			// segun args0 definir que hacer
			// comando = args0
			// controlar que dato (dni o nro de lista) sea un int
			switch comando {

			case "ingresar":
				ingresarDNI(fila, padron, dato)

			case "votar":
				votar(fila, datos)

			case "deshacer":
				deshacerVoto(fila)

			case "fin-voto":
				finalizarVoto(fila, partidos, impugnados, candidaturas)

			}
		}

		// iterar impugnados con un contador para saber la cantidad
		// imprimirResultados()
	}
}

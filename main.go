package main

import (
	"bufio"
	"fmt"
	TDACola "main/cola"
	errores "main/errores"
	TDALista "main/lista"
	"main/votos"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	MIN_DNI       = 1000000
	MAX_DNI       = 100000000
	INIT_PADRON   = 100
	INIT_PARTIDOS = 10
	POS_INVALIDA  = -1
)

//  ############### DESHACER ------------------------------------------------------------------------------------------
// func deshacerVoto(fila TDACola.Cola[votos.Votante]) {
// 	errDeshacer := fila.VerPrimero().Deshacer()
// 	if errDeshacer != nil {
// 		fmt.Fprintf(os.Stdout, "%s", errDeshacer)
// 	} else {
// 		fmt.Fprintf(os.Stdout, "OK")
// 	}
// }

// ############### INGRESAR DNI ---------------------------------------------------------------------------------------
func posicionPadron(padron []votos.Votante, dni int, inicio int, fin int) int {
	medio := (inicio + fin) / 2
	if inicio >= fin {
		return POS_INVALIDA
	}
	if padron[medio].LeerDNI() == dni {
		return medio
	}
	if padron[medio].LeerDNI() > dni {
		return posicionPadron(padron, dni, 0, medio)
	} else {
		return posicionPadron(padron, dni, medio, len(padron))
	}
}

func buscarEnPadron(padron []votos.Votante, dni int) (votos.Votante, error) {

	if posicionPadron(padron, dni, 0, len(padron)) == POS_INVALIDA {
		return nil, errores.DNIFueraPadron{}
	}
	return padron[posicionPadron(padron, dni, 0, len(padron))], nil

}

func controlarDNI(padron []votos.Votante, data []string) (votos.Votante, error) {
	//se podria controlar si len(data) > 1 pero no recuerdo si se contempla en los errores
	dni, err := strconv.Atoi(data[0])
	if err != nil || dni <= MIN_DNI || dni >= MAX_DNI {
		return nil, errores.DNIError{}
	}
	return buscarEnPadron(padron, dni)
}

func ingresarDNI(fila TDACola.Cola[votos.Votante], padron []votos.Votante, dni []string) {
	votanteIngresado, errIngresando := controlarDNI(padron, dni)
	if errIngresando == nil {
		fila.Encolar(votanteIngresado)
		fmt.Fprintf(os.Stdout, "OK")
	} else {
		fmt.Fprintf(os.Stdout, "%s", errIngresando)
	}
}

//  ############### VOTAR ----------------------------------------------------------------------------------------------
// func candidaturaValida(candidaturas []votos.TipoVoto, tipo string) bool {
// 	for i := range candidaturas {
// 		if candidaturas[i] == tipo {
// 			return true
// 		}
// 		return false
// 	}
// }

// func controlarTipo(tipo string, candidaturas []votos.TipoVoto) (votos.TipoVoto, error) {
// 	data := strings.ToUpper(tipo) // VER COMO CONVERTIR A TIPOVOTO asi ya se evalua en candidaturaValida como TipoVoto??

// 	if !candidaturaValida(candidaturas, data) {
// 		fmt.Fprintf(os.Stdout, "%s", errores.ErrorTipoVoto{})
// 		return data, errores.ErrorTipoVoto{}
// 	}

// 	return data, nil

// }

// func controlarAlt(alt string, partidos []votos.Partido) (int, error) {
// 	alternativa, errAlt := strconv.Atoi(alt)
// 	if errAlt != nil || alternativa > len(partidos) {
// 		fmt.Fprintf(os.Stdout, "%s", errores.ErrorAlternativaInvalida{})
// 		return -1, errAlt
// 	}
// 	return alternativa, errAlt
// }

// func votar(fila TDACola.Cola[votos.Votante], datos []string, candidaturas []votos.TipoVoto, partidos []votos.Partido) {
// 	if fila.EstaVacia() {
// 		fmt.Fprintf(os.Stdout, "%s", errores.FilaVacia{})
// 	} else {
// 		tipo, errTipo := controlarTipo(datos[0], candidaturas)
// 		alt, errAlt := controlarAlt(datos[1], partidos)
// 		if errAlt == nil && errTipo == nil {
// 			fila.VerPrimero().Votar(tipo, alt)
// 		}
// 	}
// }

// ############### FIN-VOTO  ------------------------------------------------------------------------------------------
func sumarVoto(voto votos.Voto, partidos []votos.Partido, candidaturas []votos.TipoVoto) {
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

// ############### Lectura Archivos de Inicio -------------------------------------------------------------------------

func prepararLista(lista *[]votos.Partido, archivoLista string) {

	(*lista)[0] = votos.CrearVotosEnBlanco()
	archivo, err := os.Open(archivoLista)
	if err != nil {
		fmt.Println(errores.ErrorLeerArchivo.Error)
	}
	defer archivo.Close()

	s := bufio.NewScanner(archivo)
	for s.Scan() {
		dividirLinea := strings.Split(s.Text(), ",")
		partido := votos.CrearPartido(dividirLinea[0], dividirLinea[1:])
		*lista = append(*lista, partido)
	}

	err = s.Err()
	if err != nil {
		fmt.Println(err)
	}

}

func leerPadron(archivoPadron string) []int {

	var temp []int
	archivo, err := os.Open(archivoPadron)
	if err != nil {
		fmt.Println(errores.ErrorLeerArchivo.Error)
	}
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

func prepararPadron(padron *[]votos.Votante, archivoPadron string) {
	// Ordenar padron en un array para despues hacer busqueda binaria (en el caso del padron)
	// Ver si podemos no leer el padron 2 veces (una hace el array simple y la otra lo hace con el struct entero pero ya ordenado)
	//y ordenar directamente el struct
	// vi que hay un par de maneras con sort.Slice pero no puedo probar que simplemente funcione asi que no me meti ahi
	temp := leerPadron(archivoPadron)
	sort.Ints(temp)
	for i := range temp {
		*padron = append(*padron, votos.CrearVotante(temp[i]))
	}
}

func prepararMesa(archivoLista, archivoPadron string) ([]votos.Partido, []votos.Votante) {
	// estructuras que vamos a usar, puse los valores de las const como placeholder pero habria que ver cuantos partidos/dni
	//trae cada archivo de prueba y ahi hacer el array? porque en caso de un archivo de 300mil va a redimensionar banda
	// no se que conviene, sobretodo en el padron, la lista de partidos suele ser corta
	padron := make([]votos.Votante, INIT_PADRON)
	lista := make([]votos.Partido, INIT_PARTIDOS)
	// leer archivos
	prepararPadron(&padron, archivoPadron)
	prepararLista(&lista, archivoLista)
	//fmt.Println(padron)
	//fmt.Println(lista)
	return lista, padron
}

func inicializar(args []string) bool {
	// tecnicamente estos mismos errores se pueden manejar con el scanner pero queria que lo comprobara antes de
	// inicializar el resto del programa

	// parametros correctos
	if len(args) < 2 {
		fmt.Fprintf(os.Stdout, "%s", errores.ErrorParametros{})
		return false
	}

	// archivos existen
	_, err1 := os.Stat(args[0])
	_, err2 := os.Stat(args[1])
	if err2 != nil || err1 != nil {
		fmt.Fprintf(os.Stdout, "%s /n", errores.ErrorLeerArchivo{})
		return false
	}
	return true
}

// ############### ---------------------------------------------------------------------------------------------------

func main() {
	var (
		padron   []votos.Votante
		partidos []votos.Partido
		//candidaturas = []votos.TipoVoto{votos.PRESIDENTE, votos.GOBERNADOR, votos.GOBERNADOR}
	//impugnados   = TDALista.CrearListaEnlazada[votos.Voto]()
	// iba a hacer un array para impugnados y dar como resultado el len del array
	// pero siento que iterar la lista al final va a ser menos costoso que redimensionar tantas veces?
	)

	argumentos := os.Args

	if inicializar(argumentos[1:]) {
		partidos, padron = prepararMesa(argumentos[1], argumentos[2])
		// cola de votantes
		fila := TDACola.CrearColaEnlazada[votos.Votante]()
		partidos[0].ObtenerResultado(1)
		// // lectura stdin
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			args := strings.Split(s.Text(), " ")
			switch args[0] {
			case "ingresar":
				ingresarDNI(&fila, padron, args[1:])

			case "votar":
				//votar(fila, args[1:], candidaturas, partidos)

				// 		// case "deshacer":
				// 		// 	deshacerVoto(fila)

				// 		// case "fin-voto":
				// 		// 	finalizarVoto(fila, partidos, impugnados, candidaturas)

			}
		}

	}
	//iterar impugnados con un contador para saber la cantidad
	//imprimirResultados()

}

package main

import (
	"bufio"
	"fmt"
	TDACola "main/cola"
	errores "main/errores"
	"main/votos"
	"os"
	"sort"
	"strconv"
	"strings"
	"tp1/votos"
)

const (
	MIN_DNI       = 1000000
	MAX_DNI       = 100000000
	INIT_PADRON   = 100
	INIT_PARTIDOS = 10
	POS_INVALIDA  = -1
)

// ############### DESHACER ------------------------------------------------------------------------------------------
func deshacerVoto(fila TDACola.Cola[votos.Votante]) {
	if fila.EstaVacia() {
		fmt.Fprintf(os.Stdout, "%s \n", errores.FilaVacia{})
	}
	errDeshacer := fila.VerPrimero().Deshacer()
	if errDeshacer != nil {
		fmt.Fprintf(os.Stdout, "%s \n", errDeshacer)
	} else {
		fmt.Fprintf(os.Stdout, "OK \n")
	}
}

// ############### INGRESAR DNI ---------------------------------------------------------------------------------------

func buscarEnPadron(padron []votos.Votante, dni int) (votos.Votante, error) {
	medio := len(padron) / 2
	if padron[medio].LeerDNI() == dni {
		return padron[medio], nil
	}
	if len(padron) <= 1 {
		return nil, errores.DNIFueraPadron{}
	}
	if padron[medio].LeerDNI() > dni {
		return buscarEnPadron(padron[:medio], dni)
	} else {
		return buscarEnPadron(padron[medio:], dni)
	}
}

func controlarDNI(padron []votos.Votante, data []string) (votos.Votante, error) {
	//Se podria controlar si len(data) > 1 pero no recuerdo si se contempla en los errores
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
		fmt.Fprintf(os.Stdout, "OK \n")
	} else {
		fmt.Fprintf(os.Stdout, "%s \n", errIngresando)
	}
}

// ############### VOTAR ----------------------------------------------------------------------------------------------
func candidaturaValida(candidaturas []votos.TipoVoto, tipo votos.TipoVoto) bool {
	for i := range candidaturas {
		if candidaturas[i] == tipo {
			return true
		}
	}
	return false
}

func pasarStringATipoVoto(tipo string) votos.TipoVoto {
	switch tipo {
	case "Presidente":
		return votos.PRESIDENTE
	case "Gobernador":
		return votos.GOBERNADOR
	case "Intendente":
		return votos.INTENDENTE
	default:
		return POS_INVALIDA
	}
}

func controlarTipo(tipo string, candidaturas []votos.TipoVoto) (votos.TipoVoto, error) {

	data := pasarStringATipoVoto(tipo)
	if !candidaturaValida(candidaturas, data) {
		fmt.Fprintf(os.Stdout, "%s \n", errores.ErrorTipoVoto{})
		return data, errores.ErrorTipoVoto{}
	}
	return data, nil

}

func controlarAlt(alt string, partidos []votos.Partido) (int, error) {
	alternativa, errAlt := strconv.Atoi(alt)
	if errAlt != nil || alternativa > len(partidos) || alternativa < 0 {
		fmt.Fprintf(os.Stdout, "%s \n", errores.ErrorAlternativaInvalida{})
		return -1, errores.ErrorAlternativaInvalida{}
	}
	return alternativa, errAlt
}

func votar(fila TDACola.Cola[votos.Votante], datos []string, candidaturas []votos.TipoVoto, partidos []votos.Partido) {
	if fila.EstaVacia() {
		fmt.Fprintf(os.Stdout, "%s \n", errores.FilaVacia{})
	} else if len(datos) != 2 {
		//Creo que esta condicion no es necesaria, porque no me parecio que la pidan en ningun lado, pero la puse
		//solamente para que no tire panic si te falto poner uno de los dos argumentos al votar
		// >>>> Esta bien que la pongas, habria que poner la misma en ingresar dni
		fmt.Fprintf(os.Stdout, "%s \n%s", errores.ErrorAlternativaInvalida{}, errores.ErrorTipoVoto{})
	} else {
		tipo, errTipo := controlarTipo(datos[0], candidaturas)
		alt, errAlt := controlarAlt(datos[1], partidos)

		if errAlt == nil && errTipo == nil {
			fila.VerPrimero().Votar(tipo, alt)
			fmt.Fprintf(os.Stdout, "OK \n")
		}
	}
}

// ############### FIN-VOTO  ------------------------------------------------------------------------------------------
func sumarVoto(voto votos.Voto, partidos []votos.Partido, candidaturas []votos.TipoVoto) {
	for i := range candidaturas {
		partidos[voto.VotoPorTipo[i]].VotadoPara(candidaturas[i])
	}
}

// No probe si funciona todavia
// >>> es necesaria una func para votos en blanco? es la posicion 0 en el array de partidos, asi que si eligieron el 0
// se suma al array
func VotosEnBlanco(votanteTerminado votos.Votante, partidos []votos.Partido, voto *votos.Voto) {
	if !voto.Impugnado {
		for i := votos.PRESIDENTE; i < votos.CANT_VOTACION; i++ {
			if voto.VotoPorTipo[i] == 0 {
				partidos[0].VotadoPara(i)
			}
		}
	}
}

// Por ahora solo funciona si no votas a las 3 candidaturas con un solo votante,
// si lo haces con 3 te tira un index out of range. Le faltaria tener en cuenta los votos en blanco
func finalizarVoto(fila TDACola.Cola[votos.Votante], partidos []votos.Partido, cantImpugnados *int, candidaturas []votos.TipoVoto) {
	voto, errFinalizar := fila.VerPrimero().FinVoto()
	if errFinalizar != nil {
		fmt.Fprintf(os.Stdout, "%s", errFinalizar)
		*cantImpugnados++
	} else {
		sumarVoto(voto, partidos, candidaturas)
		fmt.Fprintf(os.Stdout, "OK \n")
	}
	fila.Desencolar()
}

// ############### Lectura Archivos de Inicio -------------------------------------------------------------------------

func prepararLista(lista *[]votos.Partido, archivoLista string) {

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
	(*lista)[0] = votos.CrearVotosEnBlanco()

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
	//sort.Slice(*padron, func(i, j int) bool { return (*padron)[i].LeerDNI() < (*padron)[j].LeerDNI() })
	//intente hacer el sort con el slice pero cuando ejecute el programa me tiro panic asi que lo deje como estaba
	for i := range temp {
		*padron = append(*padron, votos.CrearVotante(temp[i]))
	}
}

func prepararMesa(archivoLista, archivoPadron string) ([]votos.Partido, []votos.Votante) {
	// estructuras que vamos a usar, puse los valores de las const como placeholder pero habria que ver cuantos partidos/dni
	//trae cada archivo de prueba y ahi hacer el array? porque en caso de un archivo de 300mil va a redimensionar banda
	// no se que conviene, sobretodo en el padron, la lista de partidos suele ser corta

	//las deje por defecto porque sino me tiraba un error, aparte como hacemos append, queda el arreglo
	//con muchos nil adelante y a lo ultimo los partidos y dnis, con lo cual si quisieras buscar en la lista seria
	//un problema porque la posicion 1 no tendria nada, creo..
	var padron []votos.Votante
	var lista []votos.Partido
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
		fmt.Fprintf(os.Stdout, "%s \n", errores.ErrorParametros{})
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

// Impresion de resultados -------------------------------------------------------------------------------------------

func imprimirResultados(partidos []votos.Partido, candidaturas []votos.TipoVoto) {
	for i := range candidaturas {
		fmt.Fprintf(os.Stdout, "%s: /n", candidaturas[i])
		// se podria cambiar struct de partido en blanco a que tenga nombre Votos en Blanco?
		partidos[0].ObtenerResultado()
		for j := range partidos {
			partidos[j].ObtenerResultado()
		}
		fmt.Fprintf(os.Stdout, "/n")
	}
}

func imprimirImpugnados(cantImpugnados int) {
	if cantImpugnados != 1 {
		fmt.Fprintf(os.Stdout, "Votos impugnados: %s votos/n", cantImpugnados)
	} else {
		fmt.Fprintf(os.Stdout, "Votos impugnados: %s voto/n", cantImpugnados)
	}
}

func cierreComicios(fila TDACola.Cola[votos.Votante], partidos []votos.Partido, candidaturas []votos.TipoVoto, cantImpugnados int) {

	if !fila.EstaVacia() {
		fmt.Fprintf(os.Stdout, "%s: /n", errores.ErrorCiudadanosSinVotar)
	}

	imprimirResultados(partidos, candidaturas)
	imprimirImpugnados(cantImpugnados)

}

// ############### ---------------------------------------------------------------------------------------------------

// Estoy casi seguro que los comando de ingresar, votar y deshacer funcionan bien, faltaria terminar el de fin-voto
// e imprimir todos los votos en la salida.
func main() {
	var (
		padron         []votos.Votante
		partidos       []votos.Partido
		candidaturas   = []votos.TipoVoto{votos.PRESIDENTE, votos.GOBERNADOR, votos.INTENDENTE}
		cantImpugnados int
		fila           = TDACola.CrearColaEnlazada[votos.Votante]()
	)

	argumentos := os.Args

	if inicializar(argumentos[1:]) {
		partidos, padron = prepararMesa(argumentos[1], argumentos[2])
		// cola de votantes

		partidos[0].ObtenerResultado(1)
		// // lectura stdin
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			args := strings.Split(s.Text(), " ")
			switch args[0] {
			case "ingresar":
				ingresarDNI(fila, padron, args[1:])

			case "votar":
				votar(fila, args[1:], candidaturas, partidos)

			case "deshacer":
				deshacerVoto(fila)

			case "fin-voto":
				finalizarVoto(fila, partidos, &cantImpugnados, candidaturas)

			}
		}

	}

	cierreComicios(fila, partidos, candidaturas, cantImpugnados)

}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	TDACola "tp1/cola"
	"tp1/errores"
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

	dni, err := strconv.Atoi(data[0])
	if err != nil || len(data) != 1 || dni <= MIN_DNI || dni >= MAX_DNI {
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
	if errAlt != nil || alternativa >= len(partidos) || alternativa < 0 {
		fmt.Fprintf(os.Stdout, "%s \n", errores.ErrorAlternativaInvalida{})
		return -1, errores.ErrorAlternativaInvalida{}
	}
	return alternativa, errAlt
}

func votar(fila TDACola.Cola[votos.Votante], datos []string, candidaturas []votos.TipoVoto, partidos []votos.Partido) {
	if fila.EstaVacia() {
		fmt.Fprintf(os.Stdout, "%s \n", errores.FilaVacia{})
	} else if len(datos) != 2 {
		//No es una condición contemplado en consigna, pero es necesaria para el buen funcionamiento
		//De la misma manera, si los datos no son 2, no hay forma de que el voto sea válido
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
		if voto.VotoPorTipo[i] == 0 {
			partidos[len(partidos)-1].VotadoPara(candidaturas[i])
		} else {
			partidos[voto.VotoPorTipo[i]].VotadoPara(candidaturas[i])
		}
	}
}

//func VotosEnBlanco(votanteTerminado votos.Votante, partidos []votos.Partido, voto *votos.Voto) {
//	if !voto.Impugnado {
//		for i := votos.PRESIDENTE; i < votos.CANT_VOTACION; i++ {
//			if voto.VotoPorTipo[i] == 0 {
//				partidos[0].VotadoPara(i)
//			}
//		}
//	}
//}

// Por ahora solo funciona si no votas a las 3 candidaturas con un solo votante,
// si lo haces con 3 te tira un index out of range. Le faltaria tener en cuenta los votos en blanco
func finalizarVoto(fila TDACola.Cola[votos.Votante], partidos []votos.Partido, cantImpugnados *int, candidaturas []votos.TipoVoto) {
	voto, errFinalizar := fila.VerPrimero().FinVoto()
	if errFinalizar != nil {
		fmt.Fprintf(os.Stdout, "%s", errFinalizar)
	} else {
		if voto.Impugnado {
			partidos[0].VotadoPara(votos.PRESIDENTE) // elegi presidente arbitrariamente para guardar los impugnados
		} else {
			sumarVoto(voto, partidos, candidaturas)
		}
		fmt.Fprintf(os.Stdout, "OK \n")
	}
	fila.Desencolar()
}

// ############### Lectura Archivos de Inicio -------------------------------------------------------------------------

func prepararLista(archivoLista string) []votos.Partido {
	lista := make([]votos.Partido, 0, INIT_PARTIDOS)
	archivo, err := os.Open(archivoLista)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s", errores.ErrorLeerArchivo{})
	}
	defer archivo.Close()

	s := bufio.NewScanner(archivo)
	for s.Scan() {
		dividirLinea := strings.Split(s.Text(), ",")
		partido := votos.CrearPartido(dividirLinea[0], dividirLinea[1:])
		lista = append(lista, partido)
	}
	lista[0] = votos.CrearVotosEnBlanco("Votos impugnados")
	lista = append(lista, votos.CrearVotosEnBlanco("Votos en Blanco"))

	err = s.Err()
	if err != nil {
		fmt.Println(err)
	}
	return lista
}

func leerPadron(archivoPadron string) []int {

	temp := make([]int, 0, INIT_PADRON)
	archivo, err := os.Open(archivoPadron)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s", errores.ErrorLeerArchivo{})
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

func prepararPadron(archivoPadron string) []votos.Votante {

	temp := leerPadron(archivoPadron)
	//Ordenamos el padron para poder usar búsqueda binaria
	sort.Ints(temp)
	padron := make([]votos.Votante, 0, len(temp))
	for i := range temp {
		padron = append(padron, votos.CrearVotante(temp[i]))
	}
	return padron

	//Elegimos un array en vez de lista enlazada para poder realizar la busqueda de los dni con busqueda binaria.
	//La clara desventaja es la redimension que tenga que hacer en caso de padrones muy grandes.
	//Idealmente, sabriamos la cantidad de lineas para asi poder crear el array tan grande como sea necesario.
	//PREGUNTA: ¿Podriamos para esto hacer algo de las siguientes opciones?
	//1) Leer una vez el archivo para contar las lineas antes de volver a leerlo para extraer los datos? (Habian dicho
	//que no era recomendable leerlo más de una vez)
	//2) Estimar el número de lineas usando la información provista por os.Stat() (file size) y que vamos a manejarnos
	//con datos de DNI, o sea, integers en un rango especifico?
}

func prepararMesa(archivoLista, archivoPadron string) ([]votos.Partido, []votos.Votante) {

	// leer archivos
	padron := prepararPadron(archivoPadron)
	lista := prepararLista(archivoLista)
	return lista, padron
}

func inicializar(args []string) bool {
	// Nota: Tecnicamente estos mismos errores se pueden manejar con el scanner pero queriamos que lo comprobara
	// antes de inicializar el resto del programa

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
		partidos[len(partidos)].ObtenerResultado(candidaturas[i])
		// se podria cambiar struct de partido en blanco a que tenga nombre Votos en Blanco?
		for j := 1; j < (len(partidos) - 1); j++ {
			partidos[j].ObtenerResultado(candidaturas[i])
		}
		fmt.Fprintf(os.Stdout, "/n")
	}
	partidos[0].ObtenerResultado(candidaturas[0])
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
		fmt.Fprintf(os.Stdout, "%s: /n", errores.ErrorCiudadanosSinVotar{})
	}

	imprimirResultados(partidos, candidaturas)
	imprimirImpugnados(cantImpugnados)

}

// ############### ---------------------------------------------------------------------------------------------------

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

		// lectura stdin
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

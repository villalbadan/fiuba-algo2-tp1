package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	TDACola "rerepolez/cola"
	errores "rerepolez/errores"
	votos "rerepolez/votos"
	"sort"
	"strconv"
	"strings"
)

const (
	MIN_DNI             = 0
	MAX_DNI             = 100000000
	INIT_PADRON         = 100
	INIT_PARTIDOS       = 10
	POS_INVALIDA        = -1
	ARCHIVOS_INICIO     = 2
	COMANDO             = 0 // posicion del comando a realizar segun la estructura de input: COMANDO datosDeLaAccion
	CANT_DATOS_INGRESAR = 1 // cant de datos requeridas para el comando INGRESAR DNI
	CANT_DATOS_VOTAR    = 2 // cant de datos requeridos para operacion votar
	VOTO_EN_BLANCO      = 0
	LISTA_IMPUGNADOS    = 0
	ALT_INVALIDA        = -1

	PRESIDENTE_STR  = "Presidente"
	GOBERNADOR_STR  = "Gobernador"
	INTENDENTE_STR  = "Intendente"
	TIPO_VOTO_VACIO = " "
)

//  ############### DESHACER ------------------------------------------------------------------------------------------
func deshacerVoto(fila TDACola.Cola[votos.Votante]) (int, error) {
	if fila.EstaVacia() {
		return fmt.Fprintf(os.Stdout, "%s\n", errores.FilaVacia{})
	}
	errDeshacer := fila.VerPrimero().Deshacer()

	if errDeshacer != nil {
		if errors.Is(errDeshacer, errores.ErrorVotanteFraudulento{Dni: fila.VerPrimero().LeerDNI()}) {
			fila.Desencolar()
		}
		return fmt.Fprintf(os.Stdout, "%s\n", errDeshacer)
	}
	return fmt.Fprintf(os.Stdout, "OK\n")

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
	if err != nil || len(data) != CANT_DATOS_INGRESAR || dni <= MIN_DNI || dni >= MAX_DNI {
		return nil, errores.DNIError{}
	}

	return buscarEnPadron(padron, dni)
}

func ingresarDNI(fila TDACola.Cola[votos.Votante], padron []votos.Votante, dni []string) {
	votanteIngresado, errIngresando := controlarDNI(padron, dni)
	if errIngresando == nil {
		fila.Encolar(votanteIngresado)
		fmt.Fprintf(os.Stdout, "OK\n")
	} else {
		fmt.Fprintf(os.Stdout, "%s\n", errIngresando)
	}
}

//  ############### VOTAR ----------------------------------------------------------------------------------------------
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
	case PRESIDENTE_STR:
		return votos.PRESIDENTE

	case GOBERNADOR_STR:
		return votos.GOBERNADOR

	case INTENDENTE_STR:
		return votos.INTENDENTE

	default:
		return POS_INVALIDA
	}
}

func controlarTipo(tipo string, candidaturas []votos.TipoVoto) (votos.TipoVoto, error) {

	data := pasarStringATipoVoto(tipo)
	if !candidaturaValida(candidaturas, data) {
		fmt.Fprintf(os.Stdout, "%s\n", errores.ErrorTipoVoto{})
		return data, errores.ErrorTipoVoto{}
	}
	return data, nil

}

func controlarAlt(alt string, partidos []votos.Partido) (int, error) {
	alternativa, errAlt := strconv.Atoi(alt)
	if errAlt != nil || alternativa >= len(partidos)-1 || alternativa < 0 {
		fmt.Fprintf(os.Stdout, "%s\n", errores.ErrorAlternativaInvalida{})
		return ALT_INVALIDA, errores.ErrorAlternativaInvalida{}
	}
	return alternativa, errAlt
}

func votar(fila TDACola.Cola[votos.Votante], datos []string, candidaturas []votos.TipoVoto, partidos []votos.Partido) (int, error) {
	if fila.EstaVacia() {
		return fmt.Fprintf(os.Stdout, "%s\n", errores.FilaVacia{})
	}
	if len(datos) != CANT_DATOS_VOTAR {
		//No es una condición contemplada en la consigna, pero es necesaria para el buen funcionamiento
		//De la misma manera, si los datos no son 2, no hay forma de que el voto sea válido
		return fmt.Fprintf(os.Stdout, "%s\n%s", errores.ErrorAlternativaInvalida{}, errores.ErrorTipoVoto{})
	}
	tipo, errTipo := controlarTipo(datos[0], candidaturas)
	alt, errAlt := controlarAlt(datos[1], partidos)

	if errAlt == nil && errTipo == nil {
		err := fila.VerPrimero().Votar(tipo, alt)
		if err != nil {
			fila.Desencolar()
			return fmt.Fprintf(os.Stdout, "%s\n", err)
		}
		return fmt.Fprintf(os.Stdout, "OK\n")
	}

	if errTipo != nil {
		return int(tipo), errTipo
	}
	return alt, errAlt

}

// ############### FIN-VOTO  ------------------------------------------------------------------------------------------
func sumarVoto(voto votos.Voto, partidos []votos.Partido, candidaturas []votos.TipoVoto) {
	for i := range candidaturas {
		if voto.VotoPorTipo[i] == VOTO_EN_BLANCO {
			partidos[len(partidos)-1].VotadoPara(candidaturas[i])
		} else {
			partidos[voto.VotoPorTipo[i]].VotadoPara(candidaturas[i])
		}
	}
}

func finalizarVoto(fila TDACola.Cola[votos.Votante], partidos []votos.Partido, candidaturas []votos.TipoVoto) (int, error) {
	if fila.EstaVacia() {
		return fmt.Fprintf(os.Stdout, "%s\n", errores.FilaVacia{})
	}

	voto, errFinalizar := fila.VerPrimero().FinVoto()

	if errFinalizar != nil {
		return fmt.Fprintf(os.Stdout, "%s\n", errFinalizar)
	}
	if voto.Impugnado {
		partidos[LISTA_IMPUGNADOS].VotadoPara(votos.PRESIDENTE) // elegi presidente arbitrariamente para guardar los impugnados
	} else {
		sumarVoto(voto, partidos, candidaturas)
	}

	fila.Desencolar()
	return fmt.Fprintf(os.Stdout, "OK\n")

}

// ############### Lectura Archivos de Inicio -------------------------------------------------------------------------

func prepararPartidos(archivoPartidos string) []votos.Partido {
	arregloPartidos := make([]votos.Partido, 1, INIT_PARTIDOS)
	archivo, err := os.Open(archivoPartidos)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", errores.ErrorLeerArchivo{})
	}
	defer archivo.Close()

	arregloPartidos[0] = votos.CrearVotosEnBlanco("Votos Impugnados")
	s := bufio.NewScanner(archivo)
	for s.Scan() {
		dividirLinea := strings.Split(s.Text(), ",")
		partido := votos.CrearPartido(dividirLinea[0], dividirLinea[1:])
		arregloPartidos = append(arregloPartidos, partido)
	}
	arregloPartidos = append(arregloPartidos, votos.CrearVotosEnBlanco("Votos en Blanco"))

	err = s.Err()
	if err != nil {
		fmt.Println(err)
	}
	return arregloPartidos
}

func leerPadron(archivoPadron string) []int {

	temp := make([]int, 0, INIT_PADRON)
	archivo, err := os.Open(archivoPadron)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", errores.ErrorLeerArchivo{})
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
}

func prepararMesa(archivoPartidos, archivoPadron string) ([]votos.Partido, []votos.Votante) {

	// leer archivos
	padron := prepararPadron(archivoPadron)
	partidos := prepararPartidos(archivoPartidos)
	return partidos, padron
}

func inicializar(args []string) bool {
	// Nota: Tecnicamente estos mismos errores se pueden manejar con el scanner pero queriamos que lo comprobara
	// antes de inicializar el resto del programa

	// parametros correctos
	if len(args) < ARCHIVOS_INICIO {
		fmt.Fprintf(os.Stdout, "%s\n", errores.ErrorParametros{})
		return false
	}

	// archivos existen
	_, err1 := os.Stat(args[0])
	_, err2 := os.Stat(args[1])
	if err2 != nil || err1 != nil {
		fmt.Fprintf(os.Stdout, "%s\n", errores.ErrorLeerArchivo{})
		return false
	}
	return true
}

// Impresion de resultados -------------------------------------------------------------------------------------------
func pasarTipoVotoAString(candidatura votos.TipoVoto) string {
	switch candidatura {

	case votos.PRESIDENTE:
		return PRESIDENTE_STR

	case votos.GOBERNADOR:
		return GOBERNADOR_STR

	case votos.INTENDENTE:
		return INTENDENTE_STR

	}
	return TIPO_VOTO_VACIO
}

func imprimirResultados(partidos []votos.Partido, candidaturas []votos.TipoVoto) {
	for i := range candidaturas {
		fmt.Fprintf(os.Stdout, "%s:\n", pasarTipoVotoAString(candidaturas[i]))
		// imprime votos en blanco
		fmt.Fprintln(os.Stdout, partidos[len(partidos)-1].ObtenerResultado(candidaturas[i]))
		// imprime votos de los partidos
		for j := 1; j < (len(partidos) - 1); j++ {
			fmt.Fprintln(os.Stdout, partidos[j].ObtenerResultado(candidaturas[i]))

		}
		fmt.Fprintf(os.Stdout, "\n")
	}
	//imprime impugnados
	fmt.Fprintln(os.Stdout, partidos[0].ObtenerResultado(candidaturas[0]))
}

func cierreComicios(fila TDACola.Cola[votos.Votante], partidos []votos.Partido, candidaturas []votos.TipoVoto) {

	if !fila.EstaVacia() {
		fmt.Fprintf(os.Stdout, "%s\n", errores.ErrorCiudadanosSinVotar{})
	}

	imprimirResultados(partidos, candidaturas)

}

// ############### ---------------------------------------------------------------------------------------------------

func main() {
	var (
		padron       []votos.Votante
		partidos     []votos.Partido
		candidaturas = []votos.TipoVoto{votos.PRESIDENTE, votos.GOBERNADOR, votos.INTENDENTE}
		fila         = TDACola.CrearColaEnlazada[votos.Votante]()
	)

	argumentos := os.Args

	if inicializar(argumentos[1:]) {
		partidos, padron = prepararMesa(argumentos[1], argumentos[2])

		// lectura stdin
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			args := strings.Split(s.Text(), " ")
			switch args[COMANDO] {
			case "ingresar":
				ingresarDNI(fila, padron, args[COMANDO+1:])
			case "votar":
				votar(fila, args[COMANDO+1:], candidaturas, partidos)
			case "deshacer":
				deshacerVoto(fila)
			case "fin-votar":
				finalizarVoto(fila, partidos, candidaturas)
			}
		}
		cierreComicios(fila, partidos, candidaturas)
	}
}

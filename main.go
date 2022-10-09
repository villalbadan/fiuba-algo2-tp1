package tp1

import (
	"bufio"
	"fmt"
	"os"
	TDACola "tp1/cola"
	"tp1/errores"
	TDALista "tp1/lista"
	"tp1/votos"
)

const (
	MIN_DNI = 1000000
	MAX_DNI = 100000000
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


func prepararMesa() ([]votos.Partido, []votos.Votante) {
	// leer archivos
	// ordenar padron para despues hacer busqueda binaria (en el caso del padron)
	ordenarPadron(padron)
	// estructuras que vamos a usar, puse 10 como placeholder pero habria que ver cuantos partidos/dni
	//trae cada archivo de prueba y ahi hacer el array? porque en caso de un archivo de 300mil va a redimensionar banda
	// no se que conviene
	padron := make([]votos.Votante, 10)
	partidos := make([]votos.Partido, 10)
	return partidos, padron
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
		blanco = votos.CrearVotosEnBlanco()
		candidaturas = []votos.TipoVoto{votos.PRESIDENTE, votos.GOBERNADOR, votos.GOBERNADOR}
		impugnados = TDALista.CrearListaEnlazada[votos.Voto]()
		// iba a hacer un array para impugnados y dar como resultado el len del array
		// pero siento que iterar la lista al final va a ser menos costoso que redimensionar tantas veces?
	)

	if inicializar(os.Args[1:]) {
		partidos, padron = prepararMesa()
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

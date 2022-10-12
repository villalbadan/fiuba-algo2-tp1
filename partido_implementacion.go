package votos

import (
	"fmt"
)

type partidoImplementacion struct {
	nombre     string
	votos      [CANT_VOTACION]int
	candidatos []string
}

type partidoEnBlanco struct {
	nombre string
	votos  [CANT_VOTACION]int
}

func CrearPartido(nombre string, candidatos []string) Partido {
	partido := new(partidoImplementacion)
	partido.nombre = nombre
	partido.candidatos = candidatos
	return partido
}

func CrearVotosEnBlanco(nombre string) Partido {
	blanco := new(partidoEnBlanco)
	blanco.nombre = nombre
	return blanco
}

func (partido *partidoImplementacion) VotadoPara(tipo TipoVoto) {
	partido.votos[tipo]++
}

func (partido partidoImplementacion) ObtenerResultado(tipo TipoVoto) string {
	var plural string
	if partido.votos[tipo] != 1 {
		plural = "s"
	}
	return fmt.Sprintf("%s - %s: %d voto%s", partido.nombre, partido.candidatos[tipo], partido.votos[tipo], plural)
}

func (blanco *partidoEnBlanco) VotadoPara(tipo TipoVoto) {
	blanco.votos[tipo]++
}

func (blanco partidoEnBlanco) ObtenerResultado(tipo TipoVoto) string {
	var plural string
	if blanco.votos[tipo] != 1 {
		plural = "s"
	}
	return fmt.Sprintf("%s: %d voto%s", blanco.nombre, blanco.votos[tipo], plural)
}

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
	votos [CANT_VOTACION]int
}

func CrearPartido(nombre string, candidatos []string) Partido {
	partido := new(partidoImplementacion)
	partido.nombre = nombre
	partido.candidatos = candidatos
	return partido
}

func CrearVotosEnBlanco() Partido {
	return new(partidoEnBlanco)
}

func (partido *partidoImplementacion) VotadoPara(tipo TipoVoto) {
	partido.votos[tipo]++
}

func (partido partidoImplementacion) ObtenerResultado(tipo TipoVoto) string {
	return fmt.Sprintf("%s - %s: %d votos\n", partido.nombre, partido.candidatos[tipo], partido.votos[tipo])
}

func (blanco *partidoEnBlanco) VotadoPara(tipo TipoVoto) {
	blanco.votos[tipo]++
}

func (blanco partidoEnBlanco) ObtenerResultado(tipo TipoVoto) string {
	return fmt.Sprintf("Votos en blanco: %d votos\n", blanco.votos[tipo])
}

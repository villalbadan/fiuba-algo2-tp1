package votos

import (
	"tp1/errores"
	"tp1/lista"
)

type votanteImplementacion struct {
	DNI         int
	yaVoto      bool
	voto        Voto
	ordenDeVoto lista.Lista[voto]
}

type voto struct {
	tipo  TipoVoto
	lista int
}

func CrearVotante(dni int) Votante {
	return &votanteImplementacion{DNI: dni}
}

func (votante votanteImplementacion) LeerDNI() int {
	return votante.DNI
}

func (votante *votanteImplementacion) Votar(tipo TipoVoto, alternativa int) error {
	if votante.yaVoto {
		return errores.ErrorVotanteFraudulento{}
	}

	votante.ordenDeVoto.InsertarPrimero(voto{tipo, alternativa})
	votante.voto.VotoPorTipo[tipo] = alternativa
	return nil
}

func (votante *votanteImplementacion) Deshacer() error {
	if votante.yaVoto {
		return errores.ErrorVotanteFraudulento{}
	}
	if votante.ordenDeVoto.EstaVacia() {
		return errores.ErrorNoHayVotosAnteriores{}
	}
	borrar := votante.ordenDeVoto.BorrarPrimero()
	votoAnterior := encontrarVotoAnterior(votante.ordenDeVoto, borrar.tipo)
	votante.voto.VotoPorTipo[borrar.tipo] = votoAnterior
	return nil
}

func encontrarVotoAnterior(lista lista.Lista[voto], borrado TipoVoto) int {
	for iter := lista.Iterador(); iter.HaySiguiente(); {
		if iter.VerActual().tipo == borrado {
			return iter.VerActual().lista
		}
		iter.Siguiente()
	}
	return 0
}

func (votante *votanteImplementacion) FinVoto() (Voto, error) {
	if votante.yaVoto {
		votante.voto.Impugnado = true
		return votante.voto, errores.ErrorVotanteFraudulento{}
	}
	votante.yaVoto = true
	return votante.voto, nil
}

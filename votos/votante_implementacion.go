package votos

import (
	"tp1/errores"
	"tp1/lista"
)

type votanteImplementacion struct {
	DNI         int
	yaVoto      bool
	ordenDeVoto lista.Lista[votosIndividuales]
}

type votosIndividuales struct {
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

	votante.ordenDeVoto.InsertarPrimero(votosIndividuales{tipo, alternativa})
	return nil
}

func (votante *votanteImplementacion) Deshacer() error {
	if votante.yaVoto {
		return errores.ErrorVotanteFraudulento{}
	}
	if votante.ordenDeVoto.EstaVacia() {
		return errores.ErrorNoHayVotosAnteriores{}
	}
	votante.ordenDeVoto.BorrarPrimero()
	return nil
}

func votoFinal(lista lista.Lista[votosIndividuales], votoFinal *Voto) *Voto {
	var contador TipoVoto
	for iter := lista.Iterador(); iter.HaySiguiente(); {
		if contador == CANT_VOTACION {
			return votoFinal
		}

		if votoFinal.VotoPorTipo[iter.VerActual().tipo] == 0 {
			votoFinal.VotoPorTipo[iter.VerActual().tipo] = iter.VerActual().lista
			contador++
		}

		iter.Siguiente()
	}
	return votoFinal
}

func (votante *votanteImplementacion) FinVoto() (Voto, error) {
	voto := new(Voto)
	if votante.yaVoto {
		voto.Impugnado = true
		return *voto, errores.ErrorVotanteFraudulento{}
	}
	votante.yaVoto = true
	votoFinal(votante.ordenDeVoto, voto)
	return *voto, nil
}

package votos

import (
	"rerepolez/errores"
	"rerepolez/lista"
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
	return &votanteImplementacion{DNI: dni, yaVoto: false, ordenDeVoto: lista.CrearListaEnlazada[votosIndividuales]()}
}

func (votante votanteImplementacion) LeerDNI() int {
	return votante.DNI
}

func (votante *votanteImplementacion) Votar(tipo TipoVoto, alternativa int) error {
	if votante.yaVoto {
		return errores.ErrorVotanteFraudulento{}
	}
	//if alternativa != LISTA_IMPUGNA || votante.ordenDeVoto.Largo() > 0 {
	votante.ordenDeVoto.InsertarPrimero(votosIndividuales{tipo, alternativa})
	//}
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

func votoFinal(listaVotos lista.Lista[votosIndividuales], votoFinal *Voto) *Voto {
	var contador TipoVoto
	for iter := listaVotos.Iterador(); iter.HaySiguiente(); {
		if contador == CANT_VOTACION {
			return votoFinal
		}
		if iter.VerActual().lista == LISTA_IMPUGNA {
			votoFinal.Impugnado = true
			return votoFinal
		}
		if votoFinal.VotoPorTipo[iter.VerActual().tipo] == 0 { //voto vacio
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
		return *voto, errores.ErrorVotanteFraudulento{Dni: votante.DNI}
	}

	//if votante.ordenDeVoto.EstaVacia() {
	//	return *voto, errors.New("No existe voto en curso")
	//}

	votante.yaVoto = true
	votoFinal(votante.ordenDeVoto, voto)
	return *voto, nil
}

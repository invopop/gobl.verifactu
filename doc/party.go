package doc

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

func newDestinatario(party *org.Party) []*Destinatario {
	dest := &Destinatario{
		IDDestinatario: IDDestinatario{
			NombreRazon: party.Name,
		},
	}

	if party.TaxID != nil {
		if party.TaxID.Country == l10n.ES.Tax() {
			dest.IDDestinatario.NIF = party.TaxID.Code.String()
		} else {
			dest.IDDestinatario.IDOtro = IDOtro{
				CodigoPais: party.TaxID.Country.String(),
				IDType:     "04", // Code for foreign tax IDs L7
				ID:         party.TaxID.Code.String(),
			}
		}
	}
	return []*Destinatario{dest}
}

func newTercero(party *org.Party) *Tercero {
	t := &Tercero{
		NombreRazon: party.Name,
	}

	if party.TaxID != nil {
		if party.TaxID.Country == l10n.ES.Tax() {
			t.NIF = party.TaxID.Code.String()
		} else {
			t.IDOtro = IDOtro{
				CodigoPais: party.TaxID.Country.String(),
				IDType:     "04", // Code for foreign tax IDs L7
				ID:         party.TaxID.Code.String(),
			}
		}
	}
	return t
}

package doc

import (
	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

func newDestinatario(party *org.Party) []*Destinatario {
	dest := &Destinatario{
		IDDestinatario: &Party{
			NombreRazon: party.Name,
		},
	}

	if party.TaxID != nil {
		if party.TaxID.Country == l10n.ES.Tax() {
			dest.IDDestinatario.NIF = party.TaxID.Code.String()
		} else {
			dest.IDDestinatario.IDOtro = &IDOtro{
				CodigoPais: party.TaxID.Country.String(),
				IDType:     "04", // Code for foreign tax IDs L7
				ID:         party.TaxID.Code.String(),
			}
		}
	}
	return []*Destinatario{dest}
}

func newParty(p *org.Party) *Party {
	pty := &Party{
		NombreRazon: p.Name,
	}
	if p.TaxID != nil && p.TaxID.Code.String() != "" {
		pty.NIF = p.TaxID.Code.String()
	} else {
		if len(p.Identities) > 0 {
			for _, id := range p.Identities {
				if id.Ext != nil && id.Ext[verifactu.ExtKeyIdentity] != "" {
					pty.IDOtro = &IDOtro{
						IDType: string(id.Ext[verifactu.ExtKeyIdentity]),
						ID:     id.Code.String(),
					}
				}
			}
		}
	}
	return pty
}

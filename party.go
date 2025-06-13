package verifactu

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
)

var idTypeCodeMap = map[cbc.Key]string{
	org.IdentityKeyPassport: "03",
	org.IdentityKeyForeign:  "04",
	org.IdentityKeyResident: "05",
	org.IdentityKeyOther:    "06",
}

// newParty builds a new party, but only if there are enough tax identification details
func newParty(p *org.Party) *Party {
	if p == nil {
		return nil
	}
	pty := &Party{
		NombreRazon: p.Name,
	}
	if p.TaxID != nil && p.TaxID.Code.String() != "" && p.TaxID.Country.String() == "ES" {
		pty.NIF = p.TaxID.Code.String()
	} else {
		pty.IDOtro = otherIdentity(p)
	}
	if pty.NIF == "" && pty.IDOtro == nil {
		return nil
	}
	return pty
}

func otherIdentity(p *org.Party) *IDOtro {
	oid := new(IDOtro)
	if p.TaxID != nil && p.TaxID.Code != "" {
		oid.IDType = idTypeCodeMap[org.IdentityKeyForeign]
		oid.ID = p.TaxID.Code.String()
		if p.TaxID.Country != "" {
			oid.CodigoPais = p.TaxID.Country.String()
		}
		return oid
	}

	for _, id := range p.Identities {
		it, ok := idTypeCodeMap[id.Key]
		if !ok {
			continue
		}

		oid.IDType = it
		oid.ID = id.Code.String()
		return oid
	}
	return nil
}

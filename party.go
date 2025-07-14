package verifactu

import (
	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
)

var idTypeCodeMap = map[cbc.Key]cbc.Code{
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
	if p.TaxID != nil && !p.TaxID.Code.IsEmpty() && p.TaxID.Country.In("ES") {
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
	if p.TaxID != nil {
		oid.CodigoPais = p.TaxID.Country.String()
		if p.TaxID.Code != "" {
			oid.IDType = "02" // NIF-VAT
			oid.ID = p.TaxID.Code.String()
		} else {
			// For usage with a US company for example with no
			// tax ID.
			oid.IDType = "07" // not documented (no censado)
		}
		return oid
	}

	for _, id := range p.Identities {
		code := id.Ext.Get(verifactu.ExtKeyIdentityType)

		// Fallback to matching the key, this should no longer be needed
		if code.IsEmpty() {
			if it, ok := idTypeCodeMap[id.Key]; ok {
				code = it
			}
		}
		if code.IsEmpty() {
			// Nothing to do here as there is no useful identity type
			continue
		}

		oid.CodigoPais = id.Country.String()
		oid.IDType = code.String()
		oid.ID = id.Code.String()

		return oid
	}
	return nil
}

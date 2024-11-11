package verifactu

import (
	"errors"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
)

// NewDocument creates a new  document from the provided GOBL Envelope.
// The envelope must contain a valid Invoice.
func (c *Client) Convert(env *gobl.Envelope) (*doc.VeriFactu, error) {
	// Extract the Invoice
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, errors.New("only invoices are supported")
	}
	// Check the existing stamps, we might not need to do anything
	if hasExistingStamps(env) {
		return nil, errors.New("already has stamps")
	}
	if inv.Supplier.TaxID.Country != l10n.ES.Tax() {
		return nil, errors.New("only spanish invoices are supported")
	}

	out, err := doc.NewDocument(inv, c.CurrentTime(), c.issuerRole)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Fingerprint generates a fingerprint for the  document using the
// data provided from the previous chain data. If there was no previous
// document in the chain, the parameter should be nil. The document is updated
// in place.
func (c *Client) Fingerprint(d *doc.VeriFactu) error {
	return d.Fingerprint()
}

func hasExistingStamps(env *gobl.Envelope) bool {
	for _, stamp := range env.Head.Stamps {
		if stamp.Provider.In(verifactu.StampQR) {
			return true
		}
	}
	return false
}

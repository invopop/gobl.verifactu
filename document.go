package verifactu

import (
	"github.com/invopop/gobl"
	"github.com/invopop/gobl.verifactu/internal/doc"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
)

// Document is a wrapper around the internal TicketBAI document.
type Document struct {
	env    *gobl.Envelope
	inv    *bill.Invoice
	vf     *doc.VeriFactu
	client *Client
}

// NewDocument creates a new TicketBAI document from the provided GOBL Envelope.
// The envelope must contain a valid Invoice.
func (c *Client) NewDocument(env *gobl.Envelope) (*Document, error) {
	d := new(Document)

	// Set the client for later use
	d.client = c

	var ok bool
	d.env = env
	d.inv, ok = d.env.Extract().(*bill.Invoice)
	if !ok {
		return nil, ErrOnlyInvoices
	}

	// Check the existing stamps, we might not need to do anything
	// if d.hasExistingStamps() {
	// 	return nil, ErrAlreadyProcessed
	// }
	if d.inv.Supplier.TaxID.Country != l10n.ES.Tax() {
		return nil, ErrNotSpanish
	}

	var err error
	d.vf, err = doc.NewVeriFactu(d.inv, c.CurrentTime())
	if err != nil {
		return nil, err
	}

	return d, nil
}

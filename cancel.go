package verifactu

import (
	"errors"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
)

// GenerateCancel creates a new AnulaTicketBAI document from the provided
// GOBL Envelope.
func (c *Client) GenerateCancel(env *gobl.Envelope) (*doc.VeriFactu, error) {
	// Extract the Invoice
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, errors.New("only invoices are supported")
	}
	if inv.Supplier.TaxID.Country != l10n.ES.Tax() {
		return nil, errors.New("only spanish invoices are supported")
	}
	// Extract the time when the invoice was posted to TicketBAI gateway
	// ts, err := extractPostTime(env)
	ts, err := time.Parse("2006-01-02", inv.IssueDate.String()) // REVISAR
	if err != nil {
		return nil, err
	}

	// Create the document
	cd, err := doc.NewDocument(inv, ts, c.issuerRole, c.software, true)
	if err != nil {
		return nil, err
	}

	return cd, nil
}

// FingerprintCancel generates a fingerprint for the cancellation document using the
// data provided from the previous chain data. If there was no previous
// document in the chain, the parameter should be nil. The document is updated
// in place.
func (c *Client) FingerprintCancel(d *doc.VeriFactu, prev *doc.ChainData) error {
	return d.FingerprintCancel(prev)
}

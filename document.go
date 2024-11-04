package verifactu

import (
	"github.com/invopop/gobl"
	"github.com/invopop/gobl.verifactu/internal/doc"
	"github.com/invopop/gobl/bill"
)

// Document is a wrapper around the internal TicketBAI document.
type Document struct {
	env    *gobl.Envelope
	inv    *bill.Invoice
	vf     *doc.VeriFactu
	client *Client
}

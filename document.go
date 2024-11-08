package verifactu

import (
	"errors"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/xmldsig"
)

// NewDocument creates a new Tickeverifactu document from the provided GOBL Envelope.
// The envelope must contain a valid Invoice.
func (c *Client) Convert(env *gobl.Envelope) (*doc.Tickeverifactu, error) {
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

// ZoneFor determines the zone of the envelope.
func ZoneFor(env *gobl.Envelope) l10n.Code {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return ""
	}
	return zoneFor(inv)
}

// zoneFor determines the zone of the invoice.
func zoneFor(inv *bill.Invoice) l10n.Code {
	// Figure out the zone
	if inv == nil ||
		inv.Tax == nil ||
		inv.Tax.Ext == nil ||
		inv.Tax.Ext[verifactu.ExtKeyRegion] == "" {
		return ""
	}
	return l10n.Code(inv.Tax.Ext[verifactu.ExtKeyRegion])
}

// Fingerprint generates a fingerprint for the Tickeverifactu document using the
// data provided from the previous chain data. If there was no previous
// document in the chain, the parameter should be nil. The document is updated
// in place.
func (c *Client) Fingerprint(d *doc.Tickeverifactu, prev *doc.ChainData) error {
	soft := &doc.Software{
		License: c.software.License,
		NIF:     c.software.NIF,
		Name:    c.software.Name,
		Version: c.software.Version,
	}
	return d.Fingerprint(soft, prev)
}

// Sign is used to generate the XML DSig components of the final XML document.
// This method will also update the GOBL Envelope with the QR codes that are
// generated.
func (c *Client) Sign(d *doc.Tickeverifactu, env *gobl.Envelope) error {
	zone := ZoneFor(env)
	dID := env.Head.UUID.String()
	if err := d.Sign(dID, c.cert, c.issuerRole, zone, xmldsig.WithCurrentTime(d.IssueTimestamp)); err != nil {
		return fmt.Errorf("signing: %w", err)
	}

	// now generate the QR codes and add them to the envelope
	codes := d.QRCodes(zone)
	env.Head.AddStamp(
		&head.Stamp{
			Provider: verifactu.StampCode,
			Value:    codes.verifactuCode,
		},
	)
	env.Head.AddStamp(
		&head.Stamp{
			Provider: verifactu.StampQR,
			Value:    codes.QRCode,
		},
	)
	return nil
}

func hasExistingStamps(env *gobl.Envelope) bool {
	for _, stamp := range env.Head.Stamps {
		if stamp.Provider.In(verifactu.StampCode, verifactu.StampQR) {
			return true
		}
	}
	return false
}

// Package verifactu provides the VeriFactu client
package verifactu

import (
	"context"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/xmldsig"
)

const (
	// StampKeyHash defines the key used to store the Hash in the envelope
	// stamps.
	StampKeyHash cbc.Key = "verifactu-hash"
)

// Client provides the main interface to the VeriFactu package.
type Client struct {
	software   *Software
	env        Environment
	issuerRole IssuerRole
	rep        *Issuer
	curTime    time.Time
	cert       *xmldsig.Certificate
	conn       *connection
}

// Option is used to configure the client.
type Option func(*Client)

// WithCertificate defines the signing certificate to use when producing the
// VeriFactu document.
func WithCertificate(cert *xmldsig.Certificate) Option {
	return func(c *Client) {
		c.cert = cert
	}
}

// WithCurrentTime defines the current time to use when generating the VeriFactu
// document. Only useful for testing.
func WithCurrentTime(curTime time.Time) Option {
	return func(c *Client) {
		c.curTime = curTime
	}
}

// WithSupplierIssuer set the issuer type to supplier.
func WithSupplierIssuer() Option {
	return func(c *Client) {
		c.issuerRole = IssuerRoleSupplier
	}
}

// WithRepresentative can be used to define a legal representative
// of the entity legally obliged to issue the invoice, usually the
// supplier. This may be required when the supplier's digital
// certificate is not available.
func WithRepresentative(name, nif string) Option {
	return func(c *Client) {
		c.rep = &Issuer{
			NombreRazon: name,
			NIF:         nif,
		}
	}
}

// WithCustomerIssuer set the issuer type to customer.
func WithCustomerIssuer() Option {
	return func(c *Client) {
		c.issuerRole = IssuerRoleCustomer
	}
}

// WithThirdPartyIssuer set the issuer type to third party.
func WithThirdPartyIssuer() Option {
	return func(c *Client) {
		c.issuerRole = IssuerRoleThirdParty
	}
}

// InProduction defines the connection to use the production environment.
func InProduction() Option {
	return func(c *Client) {
		c.env = EnvironmentProduction
	}
}

// InSandbox defines the connection to use the testing environment.
func InSandbox() Option {
	return func(c *Client) {
		c.env = EnvironmentSandbox
	}
}

// New creates a new VeriFactu client with shared software and configuration
// options for creating and sending new documents.
func New(software *Software, opts ...Option) (*Client, error) {
	c := new(Client)
	c.software = software

	// Set default values that can be overwritten by the options
	c.env = EnvironmentSandbox
	c.issuerRole = IssuerRoleSupplier

	for _, opt := range opts {
		opt(c)
	}

	if c.cert == nil {
		return c, nil
	}

	if c.conn == nil {
		var err error
		c.conn, err = newConnection(c.env, c.cert)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// CurrentTime returns the current time to use when generating
// the VeriFactu document.
func (c *Client) CurrentTime() time.Time {
	if !c.curTime.IsZero() {
		return c.curTime
	}
	return time.Now()
}

// Sandbox returns true if the client is using the sandbox environment.
func (c *Client) Sandbox() bool {
	return c.env == EnvironmentSandbox
}

// RegisterInvoice prepares a new registration document from the provided invoice
// inside the GOBL envelope. It will fingerprint and update the envelope with
// the chaining hash and QR code. The Registration document can be persisted for
// sending later.
func (c *Client) RegisterInvoice(env *gobl.Envelope, prev *ChainData) (*InvoiceRegistration, error) {
	if hasExistingStamps(env) {
		return nil, ErrAlreadyProcessed
	}
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, ErrOnlyInvoices
	}
	if inv.GetRegime() != l10n.ES.Tax() {
		return nil, ErrNotSpanish
	}

	if inv.Type == bill.InvoiceTypeCreditNote {
		// In VeriFactu credit notes become "facturas rectificativas por diferencias",
		// which require negative totals.
		if err := inv.Invert(); err != nil {
			return nil, err
		}
	}

	reg, err := newInvoiceRegistration(inv, c.CurrentTime(), c.issuerRole, c.software)
	if err != nil {
		return nil, err
	}
	reg.fingerprint(prev)
	c.addRegistrationStamps(env, reg)

	return reg, nil
}

// CancelInvoice builds a cancellation message from the provided document and previous
// chain data. Note that the cancellation does not require Hash information of the last,
// invoice, and instead only requires the previous chain entry.
func (c *Client) CancelInvoice(env *gobl.Envelope, prev *ChainData) (*InvoiceCancellation, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, ErrOnlyInvoices
	}
	can := newInvoiceCancellation(inv, c.CurrentTime(), c.software)
	can.fingerprint(prev)

	return can, nil
}

// NewInvoiceRequest prepares a new invoice request for the provided supplier.
func (c *Client) NewInvoiceRequest(supplier *org.Party) (*InvoiceRequest, error) {
	if supplier == nil || supplier.TaxID == nil {
		return nil, ErrValidation.WithMessage("missing supplier or tax id")
	}
	ir := new(InvoiceRequest)
	ir.Header = &InvoiceRequestHeader{
		Obligado: Issuer{
			NombreRazon: supplier.Name,
			NIF:         supplier.TaxID.Code.String(),
		},
		Representante: c.rep,
	}
	return ir, nil
}

// NewEnvelopeInvoiceRequest is a convenience method that prepares a new InvoiceRequest
// from the GOBL envelope in a single method call.
func (c *Client) NewEnvelopeInvoiceRequest(env *gobl.Envelope, prev *ChainData) (*InvoiceRequest, error) {
	req, err := c.RegisterInvoice(env, prev)
	if err != nil {
		return nil, err
	}
	inv := env.Extract().(*bill.Invoice)
	ir, err := c.NewInvoiceRequest(inv.Supplier)
	if err != nil {
		return nil, err
	}

	ir.AddRegistration(req)
	return ir, nil
}

// SendInvoiceRequest will prepare the final SOAP envelope with the invoice request
// data and send it the agency API.
func (c *Client) SendInvoiceRequest(ctx context.Context, ir *InvoiceRequest) (*InvoiceResponse, error) {
	if len(ir.Lines) == 0 {
		return nil, ErrValidation.WithMessage("no invoice request lines")
	}

	data, err := ir.Envelop().Bytes()
	if err != nil {
		return nil, err
	}

	out, err := c.conn.post(ctx, data)
	if err != nil {
		return nil, err
	}

	if res := out.Body.InvoiceResponse; res == nil {
		return nil, ErrConnection.WithMessage("missing response body")
	}
	return out.Body.InvoiceResponse, nil
}

// addRegistrationStamps adds the QR code stamp and Hash to the envelope.
func (c *Client) addRegistrationStamps(env *gobl.Envelope, reg *InvoiceRegistration) {
	// now generate the QR codes and add them to the envelope
	code := reg.generateURL(c.env == EnvironmentProduction)
	env.Head.AddStamp(&head.Stamp{
		Provider: verifactu.StampQR,
		Value:    code,
	})
	env.Head.AddStamp(&head.Stamp{
		Provider: StampKeyHash,
		Value:    reg.Huella,
	})
}

func hasExistingStamps(env *gobl.Envelope) bool {
	for _, stamp := range env.Head.Stamps {
		if stamp.Provider.In(verifactu.StampQR) {
			return true
		}
	}
	return false
}

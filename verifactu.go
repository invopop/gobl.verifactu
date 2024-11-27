// Package verifactu provides the VeriFactu client
package verifactu

import (
	"context"
	"errors"
	"time"

	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/gobl.verifactu/internal/gateways"
	"github.com/invopop/xmldsig"
)

// Standard error responses.
var (
	ErrNotSpanish       = newValidationError("only spanish invoices are supported")
	ErrAlreadyProcessed = newValidationError("already processed")
	ErrOnlyInvoices     = newValidationError("only invoices are supported")
)

// ValidationError is a simple wrapper around validation errors (that should not be retried) as opposed
// to server-side errors (that should be retried).
type ValidationError struct {
	err error
}

// Error implements the error interface for ClientError.
func (e *ValidationError) Error() string {
	return e.err.Error()
}

func newValidationError(text string) error {
	return &ValidationError{errors.New(text)}
}

// Client provides the main interface to the VeriFactu package.
type Client struct {
	software   *doc.Software
	env        gateways.Environment
	issuerRole doc.IssuerRole
	curTime    time.Time
	cert       *xmldsig.Certificate
	gw         *gateways.Connection
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
// document. Useful for testing.
func WithCurrentTime(curTime time.Time) Option {
	return func(c *Client) {
		c.curTime = curTime
	}
}

// New creates a new VeriFactu client with shared software and configuration
// options for creating and sending new documents.
func New(software *doc.Software, opts ...Option) (*Client, error) {
	c := new(Client)
	c.software = software

	// Set default values that can be overwritten by the options
	c.env = gateways.EnvironmentSandbox
	c.issuerRole = doc.IssuerRoleSupplier

	for _, opt := range opts {
		opt(c)
	}

	if c.cert == nil {
		return c, nil
	}

	if c.gw == nil {
		var err error
		c.gw, err = gateways.New(c.env, c.cert)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// WithSupplierIssuer set the issuer type to supplier.
func WithSupplierIssuer() Option {
	return func(c *Client) {
		c.issuerRole = doc.IssuerRoleSupplier
	}
}

// WithCustomerIssuer set the issuer type to customer.
func WithCustomerIssuer() Option {
	return func(c *Client) {
		c.issuerRole = doc.IssuerRoleCustomer
	}
}

// WithThirdPartyIssuer set the issuer type to third party.
func WithThirdPartyIssuer() Option {
	return func(c *Client) {
		c.issuerRole = doc.IssuerRoleThirdParty
	}
}

// InProduction defines the connection to use the production environment.
func InProduction() Option {
	return func(c *Client) {
		c.env = gateways.EnvironmentProduction
	}
}

// InSandbox defines the connection to use the testing environment.
func InSandbox() Option {
	return func(c *Client) {
		c.env = gateways.EnvironmentSandbox
	}
}

// Post will send the document to the VeriFactu gateway.
func (c *Client) Post(ctx context.Context, d *doc.VeriFactu) error {
	if err := c.gw.Post(ctx, *d); err != nil {
		return err
	}
	return nil
}

// CurrentTime returns the current time to use when generating
// the VeriFactu document.
func (c *Client) CurrentTime() time.Time {
	if !c.curTime.IsZero() {
		return c.curTime
	}
	return time.Now()
}

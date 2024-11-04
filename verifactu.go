package verifactu

import (
	"errors"
	"time"

	"github.com/invopop/gobl.verifactu/internal/doc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/xmldsig"
)

// Standard error responses.
var (
	ErrNotSpanish       = newValidationError("only spanish invoices are supported")
	ErrAlreadyProcessed = newValidationError("already processed")
	ErrOnlyInvoices     = newValidationError("only invoices are supported")
	ErrInvalidZone      = newValidationError("invalid zone")
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
	software *doc.Software
	// list       *gateways.List
	cert *xmldsig.Certificate
	// env        gateways.Environment
	// issuerRole doc.IssuerRole
	curTime time.Time
	zone    l10n.Code
}

// Option is used to configure the client.
type Option func(*Client)

// WithCurrentTime defines the current time to use when generating the VeriFactu
// document. Useful for testing.
func WithCurrentTime(curTime time.Time) Option {
	return func(c *Client) {
		c.curTime = curTime
	}
}

// PreviousInvoice stores the fields from the previously generated invoice
// document that are linked to in the new document.
type PreviousInvoice struct {
	Series    string
	Code      string
	IssueDate string
	Signature string
}

// New creates a new VeriFactu client with shared software and configuration
// options for creating and sending new documents.
func New(software *doc.Software, opts ...Option) (*Client, error) {
	c := new(Client)
	c.software = software

	// Set default values that can be overwritten by the options
	// c.env = gateways.EnvironmentTesting
	// c.issuerRole = doc.IssuerRoleSupplier

	// for _, opt := range opts {
	// 	opt(c)
	// }

	// // Create a new gateway list if none was created by the options
	// if c.list == nil && c.cert != nil {
	// 	list, err := gateways.New(c.env, c.cert)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("creating gateway list: %w", err)
	// 	}

	// 	c.list = list
	// }

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

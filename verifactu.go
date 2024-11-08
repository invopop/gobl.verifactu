package verifactu

import (
	"errors"
	"fmt"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.verifactu/internal/doc"
	"github.com/invopop/gobl.verifactu/internal/gateways"
	"github.com/invopop/gobl/bill"
	// "github.com/invopop/gobl/l10n"
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

type Software struct {
	NombreRazon                 string
	NIF                         string
	IdSistemaInformatico        string
	NombreSistemaInformatico    string
	NumeroInstalacion           string
	TipoUsoPosibleSoloVerifactu string
	TipoUsoPosibleMultiOT       string
	IndicadorMultiplesOT        string
	Version                     string
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
	software *Software
	// list       *gateways.List
	env        gateways.Environment
	issuerRole doc.IssuerRole
	curTime    time.Time
	// zone    l10n.Code
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
func New(software *Software, opts ...Option) (*Client, error) {
	c := new(Client)
	c.software = software

	// Set default values that can be overwritten by the options
	c.env = gateways.EnvironmentTesting
	c.issuerRole = doc.IssuerRoleSupplier

	for _, opt := range opts {
		opt(c)
	}

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

func (c *Client) NewVerifactu(env *gobl.Envelope) (*doc.VeriFactu, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}
	doc, err := doc.NewVeriFactu(inv, c.CurrentTime())
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// CurrentTime returns the current time to use when generating
// the VeriFactu document.
func (c *Client) CurrentTime() time.Time {
	if !c.curTime.IsZero() {
		return c.curTime
	}
	return time.Now()
}

// Package doc provides the VeriFactu document mappings from GOBL
package doc

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/invopop/gobl/bill"
)

// for needed for timezones
var location *time.Location

// IssuerRole defines the role of the issuer in the invoice.
type IssuerRole string

// IssuerRole constants
const (
	IssuerRoleSupplier   IssuerRole = "E"
	IssuerRoleCustomer   IssuerRole = "D"
	IssuerRoleThirdParty IssuerRole = "T"
)

const (
	// CurrentVersion is the current version of the VeriFactu document
	CurrentVersion = "1.0"
)

func init() {
	var err error
	location, err = time.LoadLocation("Europe/Madrid")
	if err != nil {
		panic(err)
	}
}

// Options defines the set of options to be used when building
// a new VeriFactu document.
type Options struct {
	// Software defines the software used to generate the document.
	Software *Software

	// IssuerRole defines the role of the issuer in the document.
	IssuerRole IssuerRole

	// Representative defines the legal representative of the entity
	// legally obliged to issue the invoice, usually the supplier.
	Representative *Obligado

	// Timestamp defines the current time to use when signing the
	// document.
	Timestamp time.Time
}

// NewInvoice creates a new VeriFactu invoice document from the provided GOBL Invoice.
func NewInvoice(inv *bill.Invoice, opts *Options) (*Envelope, error) {
	return newEnvelope(inv, opts, false)
}

// NewCancel creates a new VeriFactu cancellation document from the provided GOBL Invoice.
func NewCancel(inv *bill.Invoice, opts *Options) (*Envelope, error) {
	return newEnvelope(inv, opts, true)
}

func newEnvelope(inv *bill.Invoice, opts *Options, cancel bool) (*Envelope, error) {
	env := &Envelope{
		XMLNs: EnvNamespace,
		SUM:   SUM,
		SUM1:  SUM1,
		Body: &Body{
			VeriFactu: &RegFactuSistemaFacturacion{},
		},
	}

	doc := &RegFactuSistemaFacturacion{
		Cabecera: &Cabecera{
			Obligado: Obligado{
				NombreRazon: inv.Supplier.Name,
				NIF:         inv.Supplier.TaxID.Code.String(),
			},
		},
		RegistroFactura: &RegistroFactura{},
	}
	if opts.Representative != nil {
		doc.Cabecera.Representante = opts.Representative
	}

	if inv.Type == bill.InvoiceTypeCreditNote {
		// GOBL credit and debit notes' amounts represent the amounts to be credited to the customer,
		// and they are provided as positive numbers. In VeriFactu, however, credit notes
		// become "facturas rectificativas por diferencias" and, when a correction is for a
		// credit operation, the amounts must be negative to cancel out the ones in the
		// original invoice. For that reason, we invert the credit note quantities here.
		if err := inv.Invert(); err != nil {
			return nil, err
		}
	}

	if cancel {
		reg := newCancel(inv, opts.Timestamp, opts.Software)
		doc.RegistroFactura.RegistroAnulacion = reg
	} else {
		reg, err := newInvoice(inv, opts.Timestamp, opts.IssuerRole, opts.Software)
		if err != nil {
			return nil, err
		}
		doc.RegistroFactura.RegistroAlta = reg
	}

	env.Body.VeriFactu = doc

	return env, nil
}

// QRCodes generates the QR code for the document
func (d *Envelope) QRCodes(production bool) string {
	return d.generateURL(production)
}

// ChainData generates the data to be used to link to this one
// in the next entry.
func (d *Envelope) ChainData() Encadenamiento {
	return Encadenamiento{
		RegistroAnterior: &RegistroAnterior{
			IDEmisorFactura:        d.Body.VeriFactu.Cabecera.Obligado.NIF,
			NumSerieFactura:        d.Body.VeriFactu.RegistroFactura.RegistroAlta.IDFactura.NumSerieFactura,
			FechaExpedicionFactura: d.Body.VeriFactu.RegistroFactura.RegistroAlta.IDFactura.FechaExpedicionFactura,
			Huella:                 d.Body.VeriFactu.RegistroFactura.RegistroAlta.Huella,
		},
	}
}

// ChainDataCancel generates the data to be used to link to this one
// in the next entry for cancelling invoices.
func (d *Envelope) ChainDataCancel() Encadenamiento {
	return Encadenamiento{
		RegistroAnterior: &RegistroAnterior{
			IDEmisorFactura:        d.Body.VeriFactu.Cabecera.Obligado.NIF,
			NumSerieFactura:        d.Body.VeriFactu.RegistroFactura.RegistroAnulacion.IDFactura.NumSerieFactura,
			FechaExpedicionFactura: d.Body.VeriFactu.RegistroFactura.RegistroAnulacion.IDFactura.FechaExpedicionFactura,
			Huella:                 d.Body.VeriFactu.RegistroFactura.RegistroAnulacion.Huella,
		},
	}
}

// Fingerprint generates the SHA-256 fingerprint for the document
func (d *Envelope) Fingerprint(prev *ChainData) error {
	return d.generateHashAlta(prev)
}

// FingerprintCancel generates the SHA-256 fingerprint for the document
func (d *Envelope) FingerprintCancel(prev *ChainData) error {
	return d.generateHashAnulacion(prev)
}

// Bytes returns the XML document bytes
func (d *Envelope) Bytes() ([]byte, error) {
	return toBytes(d)
}

// Bytes returns the XML document bytes
func (d *RegFactuSistemaFacturacion) Bytes() ([]byte, error) {
	return toBytes(d)
}

// BytesIndent returns the indented XML document bytes
func (d *Envelope) BytesIndent() ([]byte, error) {
	return toBytesIndent(d)
}

func toBytes(doc any) ([]byte, error) {
	buf, err := buffer(doc, xml.Header, false)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func toBytesIndent(doc any) ([]byte, error) {
	buf, err := buffer(doc, xml.Header, true)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buffer(doc any, base string, indent bool) (*bytes.Buffer, error) {
	buf := bytes.NewBufferString(base)

	enc := xml.NewEncoder(buf)
	if indent {
		enc.Indent("", "  ")
	}

	if err := enc.Encode(doc); err != nil {
		return nil, fmt.Errorf("encoding document: %w", err)
	}

	return buf, nil
}

type timeLocationable interface {
	In(*time.Location) time.Time
}

func formatDateTimeZone(ts timeLocationable) string {
	return ts.In(location).Format("2006-01-02T15:04:05-07:00")
}

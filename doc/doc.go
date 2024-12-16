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

// NewVerifactu creates a new VeriFactu document
func NewVerifactu(inv *bill.Invoice, ts time.Time, r IssuerRole, s *Software, c bool) (*VeriFactu, error) {
	doc := &VeriFactu{
		Cabecera: &Cabecera{
			Obligado: Obligado{
				NombreRazon: inv.Supplier.Name,
				NIF:         inv.Supplier.TaxID.Code.String(),
			},
		},
		RegistroFactura: &RegistroFactura{},
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

	if c {
		reg, err := NewCancel(inv, ts, s)
		if err != nil {
			return nil, err
		}
		doc.RegistroFactura.RegistroAnulacion = reg
	} else {
		reg, err := NewInvoice(inv, ts, r, s)
		if err != nil {
			return nil, err
		}
		doc.RegistroFactura.RegistroAlta = reg
	}

	return doc, nil
}

// QRCodes generates the QR code for the document
func (d *VeriFactu) QRCodes(production bool) string {
	return d.generateURL(production)
}

// ChainData generates the data to be used to link to this one
// in the next entry.
func (d *VeriFactu) ChainData() Encadenamiento {
	return Encadenamiento{
		RegistroAnterior: &RegistroAnterior{
			IDEmisorFactura:        d.Cabecera.Obligado.NIF,
			NumSerieFactura:        d.RegistroFactura.RegistroAlta.IDFactura.NumSerieFactura,
			FechaExpedicionFactura: d.RegistroFactura.RegistroAlta.IDFactura.FechaExpedicionFactura,
			Huella:                 d.RegistroFactura.RegistroAlta.Huella,
		},
	}
}

// ChainDataCancel generates the data to be used to link to this one
// in the next entry for cancelling invoices.
func (d *VeriFactu) ChainDataCancel() Encadenamiento {
	return Encadenamiento{
		RegistroAnterior: &RegistroAnterior{
			IDEmisorFactura:        d.Cabecera.Obligado.NIF,
			NumSerieFactura:        d.RegistroFactura.RegistroAnulacion.IDFactura.NumSerieFactura,
			FechaExpedicionFactura: d.RegistroFactura.RegistroAnulacion.IDFactura.FechaExpedicionFactura,
			Huella:                 d.RegistroFactura.RegistroAnulacion.Huella,
		},
	}
}

// Fingerprint generates the SHA-256 fingerprint for the document
func (d *VeriFactu) Fingerprint(prev *ChainData) error {
	return d.generateHashAlta(prev)
}

// FingerprintCancel generates the SHA-256 fingerprint for the document
func (d *VeriFactu) FingerprintCancel(prev *ChainData) error {
	return d.generateHashAnulacion(prev)
}

// Bytes returns the XML document bytes
func (d *VeriFactu) Bytes() ([]byte, error) {
	return toBytes(d)
}

// BytesIndent returns the indented XML document bytes
func (d *VeriFactu) BytesIndent() ([]byte, error) {
	return toBytesIndent(d)
}

// Envelop wraps the VeriFactu document in a SOAP envelope and includes the expected namespaces
func (d *VeriFactu) Envelop() ([]byte, error) {
	// Create and set the envelope with namespaces
	env := Envelope{
		XMLNs: EnvNamespace,
		SUM:   SUM,
		SUM1:  SUM1,
	}
	env.Body.VeriFactu = d

	// Marshal the SOAP envelope into an XML byte slice
	var result bytes.Buffer
	enc := xml.NewEncoder(&result)
	enc.Indent("", "  ")
	err := enc.Encode(env)
	if err != nil {
		return nil, err
	}

	// Return the enveloped XML document
	return result.Bytes(), nil
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

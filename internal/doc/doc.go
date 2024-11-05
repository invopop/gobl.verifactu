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

// ValidationError is a simple wrapper around validation errors
type ValidationError struct {
	text string
}

// Error implements the error interface for ValidationError.
func (e *ValidationError) Error() string {
	return e.text
}

func validationErr(text string, args ...any) error {
	return &ValidationError{
		text: fmt.Sprintf(text, args...),
	}
}

func init() {
	var err error
	location, err = time.LoadLocation("Europe/Madrid")
	if err != nil {
		panic(err)
	}
}

func NewVeriFactu(inv *bill.Invoice, ts time.Time, role IssuerRole) (*VeriFactu, error) {

	doc := &VeriFactu{
		SUMNamespace:  SUM,
		SUM1Namespace: SUM1,
		Cabecera: &Cabecera{
			Obligado: Obligado{
				NombreRazon: inv.Supplier.Name,
				NIF:         inv.Supplier.TaxID.Code.String(),
			},
		},
		RegistroFactura: &RegistroFactura{},
	}

	if inv.Type == bill.InvoiceTypeCreditNote {
		reg, err := NewRegistroAnulacion(inv, ts, role)
		if err != nil {
			return nil, err
		}
		doc.RegistroFactura.RegistroAnulacion = reg
	} else {
		reg, err := NewRegistroAlta(inv, ts, role)
		if err != nil {
			return nil, err
		}
		doc.RegistroFactura.RegistroAlta = reg
	}

	return doc, nil
}

func (doc *VeriFactu) QRCodes() *Codes {
	if doc.RegistroFactura.RegistroAlta.Encadenamiento == nil {
		return nil
	}
	return doc.generateCodes()
}

// Bytes returns the XML document bytes
func (doc *VeriFactu) Bytes() ([]byte, error) {
	return toBytes(doc)
}

// BytesIndent returns the indented XML document bytes
func (doc *VeriFactu) BytesIndent() ([]byte, error) {
	return toBytesIndent(doc)
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

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

type VeriFactu struct {
	Cabecera        *Cabecera
	RegistroFactura *RegistroFactura
}

type RegistroFactura struct {
	RegistroAlta      *RegistroAlta
	RegistroAnulacion *RegistroAnulacion
}

func init() {
	var err error
	location, err = time.LoadLocation("Europe/Madrid")
	if err != nil {
		panic(err)
	}
}

func NewVeriFactu(inv *bill.Invoice, ts time.Time) (*VeriFactu, error) {
	if inv.Type == bill.InvoiceTypeCreditNote {

		if err := inv.Invert(); err != nil {
			return nil, err
		}
	}

	// goblWithoutIncludedTaxes, err := inv.RemoveIncludedTaxes()
	// if err != nil {
	// 	return nil, err
	// }

	doc := &VeriFactu{
		Cabecera: &Cabecera{
			Obligado: Obligado{
				NombreRazon: inv.Supplier.Name,
				NIF:         inv.Supplier.TaxID.Code.String(),
			},
		},
		RegistroFactura: &RegistroFactura{},
	}

	doc.RegistroFactura.RegistroAlta.FechaHoraHusoGenRegistro = formatDateTimeZone(ts)

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

func toBytesCanonical(doc any) ([]byte, error) {
	buf, err := buffer(doc, "", false)
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

package doc

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/invopop/gobl/bill"
)

type VeriFactu struct {
	Cabecera        *Cabecera
	RegistroFactura *RegistroFactura
}

type RegistroFactura struct {
	RegistroAlta      *RegistroAlta
	RegistroAnulacion *RegistroAnulacion
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

func NewVeriFactu(inv *bill.Invoice, ts time.Time) (*VeriFactu, error) {
	if inv.Type == bill.InvoiceTypeCreditNote {

		if err := inv.Invert(); err != nil {
			return nil, err
		}
	}

	goblWithoutIncludedTaxes, err := inv.RemoveIncludedTaxes()
	if err != nil {
		return nil, err
	}

	doc := &VeriFactu{
		Cabecera: &Cabecera{
			Obligado: Obligado{
				NombreRazon: inv.Supplier.Name,
				NIF:         inv.Supplier.TaxID.Code.String(),
			},
		},
		RegistroFactura: &RegistroFactura{},
	}

	doc.SetIssueTimestamp(ts)

	// Add customers
	if inv.Customer != nil {
		dest, err := newDestinatario(inv.Customer)
		if err != nil {
			return nil, err
		}
		doc.Sujetos.Destinatarios = &Destinatarios{
			IDDestinatario: []*IDDestinatario{dest},
		}
	}

	if inv.Type == bill.InvoiceTypeCreditNote {
		doc.RegistroFactura.RegistroAnulacion, err = newDatosFactura(goblWithoutIncludedTaxes)
	} else {
		doc.RegistroFactura.RegistroAlta, err = newDatosFactura(goblWithoutIncludedTaxes)
	}
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// QRCodes generates the QR codes for this invoice, but requires the Fingerprint to have been
// generated first.
func (doc *TicketBAI) QRCodes() *Codes {
	if doc.HuellaTBAI == nil {
		return nil
	}
	return doc.generateCodes(doc.zone)
}

// Bytes returns the XML document bytes
func (doc *TicketBAI) Bytes() ([]byte, error) {
	return toBytes(doc)
}

// BytesIndent returns the indented XML document bytes
func (doc *TicketBAI) BytesIndent() ([]byte, error) {
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

// TODO
// func formatDate(ts timeLocationable) string {
// 	return ts.In(location).Format("02-01-2006")
// }

// func formatTime(ts timeLocationable) string {
// 	return ts.In(location).Format("15:04:05")
// }

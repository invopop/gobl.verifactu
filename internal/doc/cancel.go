package doc

import (
	"time"

	"github.com/invopop/gobl/bill"
)

type RegistroAnulacion struct {
	IDVersion                string          `xml:"IDVersion"`
	IDFactura                *IDFactura      `xml:"IDFactura"`
	RefExterna               string          `xml:"RefExterna,omitempty"`
	SinRegistroPrevio        string          `xml:"SinRegistroPrevio"`
	RechazoPrevio            string          `xml:"RechazoPrevio,omitempty"`
	GeneradoPor              string          `xml:"GeneradoPor"`
	Generador                *Tercero        `xml:"Generador"`
	Encadenamiento           *Encadenamiento `xml:"Encadenamiento"`
	SistemaInformatico       *Software       `xml:"SistemaInformatico"`
	FechaHoraHusoGenRegistro string          `xml:"FechaHoraHusoGenRegistro"`
	TipoHuella               string          `xml:"TipoHuella"`
	Huella                   string          `xml:"Huella"`
	Signature                string          `xml:"Signature"`
}

// NewRegistroAnulacion provides support for credit notes
func NewRegistroAnulacion(inv *bill.Invoice, ts time.Time) (*RegistroAnulacion, error) {
	reg := &RegistroAnulacion{
		IDVersion: "1.0",
		IDFactura: &IDFactura{
			IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
			NumSerieFactura:        invoiceNumber(inv.Series, inv.Code),
			FechaExpedicionFactura: inv.IssueDate.Time().Format("02-01-2006"),
		},
		SinRegistroPrevio: "N",
		GeneradoPor:       "1", // Generated by issuer
		Generador: &Tercero{
			Nif:         inv.Supplier.TaxID.Code.String(),
			NombreRazon: inv.Supplier.Name,
		},
		FechaHoraHusoGenRegistro: formatDateTimeZone(ts),
		TipoHuella:               "01",
	}

	return reg, nil

}

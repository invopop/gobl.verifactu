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
	Generador                *Party          `xml:"Generador"`
	Encadenamiento           *Encadenamiento `xml:"Encadenamiento"`
	SistemaInformatico       *Software       `xml:"SistemaInformatico"`
	FechaHoraHusoGenRegistro string          `xml:"FechaHoraHusoGenRegistro"`
	TipoHuella               string          `xml:"TipoHuella"`
	Huella                   string          `xml:"Huella"`
	Signature                string          `xml:"Signature"`
}

// NewRegistroAnulacion provides support for credit notes
func NewRegistroAnulacion(inv *bill.Invoice, ts time.Time, r IssuerRole, s *Software) (*RegistroAnulacion, error) {
	reg := &RegistroAnulacion{
		IDVersion: CurrentVersion,
		IDFactura: &IDFactura{
			IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
			NumSerieFactura:        invoiceNumber(inv.Series, inv.Code),
			FechaExpedicionFactura: inv.IssueDate.Time().Format("02-01-2006"),
		},
		// SinRegistroPrevio: "N",
		GeneradoPor:              string(r),
		Generador:                makeGenerador(inv, r),
		SistemaInformatico:       newSoftware(s),
		FechaHoraHusoGenRegistro: formatDateTimeZone(ts),
		TipoHuella:               "01",
	}

	return reg, nil

}

func makeGenerador(inv *bill.Invoice, r IssuerRole) *Party {
	switch r {
	case IssuerRoleSupplier, IssuerRoleThirdParty:
		return newParty(inv.Supplier)
	case IssuerRoleCustomer:
		return newParty(inv.Customer)
	}
	return nil
}

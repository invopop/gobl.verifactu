package doc

import (
	"time"

	"github.com/invopop/gobl/bill"
)

// NewRegistroAnulacion provides support for credit notes
func NewRegistroAnulacion(inv *bill.Invoice, ts time.Time, r IssuerRole, s *Software) (*RegistroAnulacion, error) {
	reg := &RegistroAnulacion{
		IDVersion: CurrentVersion,
		IDFactura: &IDFactura{
			IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
			NumSerieFactura:        invoiceNumber(inv.Series, inv.Code),
			FechaExpedicionFactura: inv.IssueDate.Time().Format("02-01-2006"),
		},
		// SinRegistroPrevio: "N", // TODO: Think what to do with this field
		// RechazoPrevio:            "N", // TODO: Think what to do with this field
		GeneradoPor:              string(r),
		Generador:                makeGenerador(inv, r),
		SistemaInformatico:       s,
		FechaHoraHusoGenRegistro: formatDateTimeZone(ts),
		TipoHuella:               TipoHuella,
	}

	return reg, nil

}

func makeGenerador(inv *bill.Invoice, r IssuerRole) *Party {
	switch r {
	case IssuerRoleSupplier, IssuerRoleThirdParty:
		p, err := newParty(inv.Supplier)
		if err != nil {
			return nil
		}
		return p
	case IssuerRoleCustomer:
		p, err := newParty(inv.Customer)
		if err != nil {
			return nil
		}
		return p
	}
	return nil
}

package doc

import (
	"time"

	"github.com/invopop/gobl/bill"
)

// newCancel provides support for cancelling invoices
func newCancel(inv *bill.Invoice, ts time.Time, s *Software) *RegistroAnulacion {
	reg := &RegistroAnulacion{
		IDVersion: CurrentVersion,
		IDFactura: &IDFacturaAnulada{
			IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
			NumSerieFactura:        invoiceNumber(inv.Series, inv.Code),
			FechaExpedicionFactura: inv.IssueDate.Time().Format("02-01-2006"),
		},
		SistemaInformatico:       s,
		FechaHoraHusoGenRegistro: formatDateTimeZone(ts),
		TipoHuella:               TipoHuella,
	}
	return reg
}

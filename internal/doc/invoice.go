package doc

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
)

func newInvoice(inv *bill.Invoice) (*RegistroAlta, error) {
	// Create new RegistroAlta with required fields
	reg := &RegistroAlta{
		IDVersion: "1.0",
		IDFactura: IDFactura{
			IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
			NumSerieFactura:        invoiceNumber(inv.Series, inv.Code),
			FechaExpedicionFactura: inv.IssueDate.Time().Format("02-01-2006"),
		},
		NombreRazonEmisor:    inv.Supplier.Name,
		FechaOperacion:       inv.IssueDate.Format("02-01-2006"),
		DescripcionOperacion: inv.Notes.String(),
		ImporteTotal:         inv.Totals.Total.Float64(),
		CuotaTotal:           inv.Totals.Tax.Float64(),
	}

	// Set TipoFactura based on invoice type
	switch inv.Type {
	case bill.InvoiceTypeStandard:
		reg.TipoFactura = "F1"
	case bill.InvoiceTypeCreditNote:
		reg.TipoFactura = "R1"
		reg.TipoRectificativa = "I" // Por diferencias
	case bill.InvoiceTypeDebitNote:
		reg.TipoFactura = "R1"
		reg.TipoRectificativa = "I"
	}

	// Add destinatarios if customer exists
	if inv.Customer != nil {
		dest := &Destinatario{
			IDDestinatario: IDDestinatario{
				NombreRazon: inv.Customer.Name,
			},
		}

		// Handle tax ID
		if inv.Customer.TaxID != nil {
			if inv.Customer.TaxID.Country.Is("ES") {
				dest.IDDestinatario.NIF = inv.Customer.TaxID.Code.String()
			} else {
				dest.IDDestinatario.IDOtro = IDOtro{
					CodigoPais: inv.Customer.TaxID.Country.String(),
					IDType:     "04", // NIF-IVA
					ID:         inv.Customer.TaxID.Code.String(),
				}
			}
		}

		reg.Destinatarios = []*Destinatario{dest}
	}

	return reg, nil
}

func invoiceNumber(series cbc.Code, code cbc.Code) string {
	if series == "" {
		return code.String()
	}
	return fmt.Sprintf("%s-%s", series, code)
}

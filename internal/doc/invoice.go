package doc

import (
	"fmt"
	"time"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

func NewRegistroAlta(inv *bill.Invoice, ts time.Time) (*RegistroAlta, error) {
	description, err := newDescription(inv.Notes)
	if err != nil {
		return nil, err
	}

	desglose, err := newDesglose(inv)
	if err != nil {
		return nil, err
	}

	reg := &RegistroAlta{
		IDVersion: "1.0",
		IDFactura: &IDFactura{
			IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
			NumSerieFactura:        invoiceNumber(inv.Series, inv.Code),
			FechaExpedicionFactura: inv.IssueDate.Time().Format("02-01-2006"),
		},
		NombreRazonEmisor:        inv.Supplier.Name,
		TipoFactura:              mapInvoiceType(inv),
		DescripcionOperacion:     description,
		ImporteTotal:             newImporteTotal(inv),
		CuotaTotal:               newTotalTaxes(inv),
		SistemaInformatico:       newSoftware(),
		Desglose:                 desglose,
		FechaHoraHusoGenRegistro: formatDateTimeZone(ts),
	}

	if inv.Customer != nil {
		dest := &Destinatario{
			IDDestinatario: IDDestinatario{
				NombreRazon: inv.Customer.Name,
			},
		}

		if inv.Customer.TaxID != nil {
			if inv.Customer.TaxID.Country == l10n.ES.Tax() {
				dest.IDDestinatario.NIF = inv.Customer.TaxID.Code.String()
			} else {
				dest.IDDestinatario.IDOtro = IDOtro{
					CodigoPais: inv.Customer.TaxID.Country.String(),
					IDType:     "04", // Code for foreign tax IDs L7
					ID:         inv.Customer.TaxID.Code.String(),
				}
			}
		}

		reg.Destinatarios = []*Destinatario{dest}
	}

	if inv.HasTags(tax.TagSimplified) {
		if inv.Type == bill.InvoiceTypeStandard {
			reg.FacturaSimplificadaArt7273 = "S"
		} else {
			reg.FacturaSinIdentifDestinatarioArt61d = "S"
		}
	}

	return reg, nil
}

func invoiceNumber(series cbc.Code, code cbc.Code) string {
	if series == "" {
		return code.String()
	}
	return fmt.Sprintf("%s-%s", series, code)
}

func mapInvoiceType(inv *bill.Invoice) string {
	switch inv.Type {
	case bill.InvoiceTypeStandard:
		return "F1"
	case bill.ShortSchemaInvoice:
		return "F2"
	}
	return "F1"
}

func newDescription(notes []*cbc.Note) (string, error) {
	for _, note := range notes {
		if note.Key == cbc.NoteKeyGeneral {
			return note.Text, nil
		}
	}
	return "", validationErr(`notes: missing note with key '%s'`, cbc.NoteKeyGeneral)
}

func newImporteTotal(inv *bill.Invoice) string {
	totalWithDiscounts := inv.Totals.Total

	totalTaxes := num.MakeAmount(0, 2)
	for _, category := range inv.Totals.Taxes.Categories {
		if !category.Retained {
			totalTaxes = totalTaxes.Add(category.Amount)
		}
	}

	return totalWithDiscounts.Add(totalTaxes).String()
}

func newTotalTaxes(inv *bill.Invoice) string {
	totalTaxes := num.MakeAmount(0, 2)
	for _, category := range inv.Totals.Taxes.Categories {
		if !category.Retained {
			totalTaxes = totalTaxes.Add(category.Amount)
		}
	}

	return totalTaxes.String()
}

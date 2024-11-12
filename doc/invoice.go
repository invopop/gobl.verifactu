package doc

import (
	"fmt"
	"time"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// NewRegistroAlta creates a new VeriFactu registration for an invoice.
func NewRegistroAlta(inv *bill.Invoice, ts time.Time, r IssuerRole, s *Software) (*RegistroAlta, error) {
	description, err := newDescription(inv.Notes)
	if err != nil {
		return nil, err
	}

	desglose := newDesglose(inv)

	reg := &RegistroAlta{
		IDVersion: CurrentVersion,
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
		SistemaInformatico:       newSoftware(s),
		Desglose:                 desglose,
		FechaHoraHusoGenRegistro: formatDateTimeZone(ts),
		TipoHuella:               TipoHuella,
	}

	if inv.Customer != nil {
		reg.Destinatarios = newDestinatario(inv.Customer)
	}

	if r == IssuerRoleThirdParty {
		reg.EmitidaPorTerceroODestinatario = "T"
		reg.Tercero = newParty(inv.Supplier)
	}

	// Check
	if inv.Type == bill.InvoiceTypeCorrective {
		reg.Subsanacion = "S"
	}

	// Check
	if inv.HasTags(tax.TagSimplified) {
		if inv.Type == bill.InvoiceTypeStandard {
			reg.FacturaSimplificadaArt7273 = "S"
		} else {
			reg.FacturaSinIdentifDestinatarioArt61d = "S"
		}
	}

	// Flag for operations with totals over 100,000,000€. Added with optimism.
	if inv.Totals.TotalWithTax.Compare(num.MakeAmount(100000000, 0)) == 1 {
		reg.Macrodato = "S"
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

func newSoftware(s *Software) *Software {
	return &Software{
		NombreRazon:          s.NombreRazon,
		NIF:                  s.NIF,
		IdSistemaInformatico: s.IdSistemaInformatico,
		Version:              s.Version,
		NumeroInstalacion:    s.NumeroInstalacion,
	}
}

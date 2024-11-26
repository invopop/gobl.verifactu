package doc

import (
	"fmt"
	"slices"
	"time"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
)

var (
	rectificative = []string{"R1", "R2", "R3", "R4", "R5", "R6"}
)

// NewRegistroAlta creates a new VeriFactu registration for an invoice.
func NewRegistroAlta(inv *bill.Invoice, ts time.Time, r IssuerRole, s *Software) (*RegistroAlta, error) {
	description, err := newDescription(inv.Notes)
	if err != nil {
		return nil, err
	}

	desglose, err := newDesglose(inv)
	if err != nil {
		return nil, err
	}

	reg := &RegistroAlta{
		IDVersion: CurrentVersion,
		IDFactura: &IDFactura{
			IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
			NumSerieFactura:        invoiceNumber(inv.Series, inv.Code),
			FechaExpedicionFactura: inv.IssueDate.Time().Format("02-01-2006"),
		},
		NombreRazonEmisor:        inv.Supplier.Name,
		TipoFactura:              inv.Tax.Ext[verifactu.ExtKeyDocType].String(),
		DescripcionOperacion:     description,
		Desglose:                 desglose,
		CuotaTotal:               newTotalTaxes(inv),
		ImporteTotal:             newImporteTotal(inv),
		SistemaInformatico:       s,
		FechaHoraHusoGenRegistro: formatDateTimeZone(ts),
		TipoHuella:               TipoHuella,
	}

	if inv.Customer != nil {
		d, err := newParty(inv.Customer)
		if err != nil {
			return nil, err
		}
		ds := &Destinatario{
			IDDestinatario: d,
		}
		reg.Destinatarios = []*Destinatario{ds}
	} else {
		reg.FacturaSinIdentifDestinatarioArt61d = "S"
	}

	if slices.Contains(rectificative, reg.TipoFactura) {
		// GOBL does not currently have explicit support for Facturas Rectificativas por Sustitución
		reg.TipoRectificativa = "I"
		if inv.Preceding != nil {
			rs := make([]*FacturaRectificada, 0, len(inv.Preceding))
			for _, ref := range inv.Preceding {
				rs = append(rs, &FacturaRectificada{
					IDFactura: IDFactura{
						IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
						NumSerieFactura:        invoiceNumber(ref.Series, ref.Code),
						FechaExpedicionFactura: ref.IssueDate.Time().Format("02-01-2006"),
					},
				})
			}
			reg.FacturasRectificadas = rs
		}
	}

	if inv.HasTags(verifactu.TagSubstitution) {
		if inv.Preceding != nil {
			subs := make([]*FacturaSustituida, 0, len(inv.Preceding))
			for _, ref := range inv.Preceding {
				subs = append(subs, &FacturaSustituida{
					IDFactura: IDFactura{
						IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
						NumSerieFactura:        invoiceNumber(ref.Series, ref.Code),
						FechaExpedicionFactura: ref.IssueDate.Time().Format("02-01-2006"),
					},
				})
			}
			reg.FacturasSustituidas = subs
		}
	}

	if r == IssuerRoleThirdParty {
		reg.EmitidaPorTerceroODestinatario = "T"
		t, err := newParty(inv.Supplier)
		if err != nil {
			return nil, err
		}
		reg.Tercero = t
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

func newDescription(notes []*cbc.Note) (string, error) {
	for _, note := range notes {
		if note.Key == cbc.NoteKeyGeneral {
			return note.Text, nil
		}
	}
	return "", validationErr(`notes: missing note with key '%s'`, cbc.NoteKeyGeneral)
}

func newImporteTotal(inv *bill.Invoice) float64 {
	totalWithDiscounts := inv.Totals.Total

	totalTaxes := num.MakeAmount(0, 2)
	for _, category := range inv.Totals.Taxes.Categories {
		if !category.Retained {
			totalTaxes = totalTaxes.Add(category.Amount)
		}
	}

	return totalWithDiscounts.Add(totalTaxes).Float64()
}

func newTotalTaxes(inv *bill.Invoice) float64 {
	totalTaxes := num.MakeAmount(0, 2)
	for _, category := range inv.Totals.Taxes.Categories {
		if !category.Retained {
			totalTaxes = totalTaxes.Add(category.Amount)
		}
	}

	return totalTaxes.Float64()
}

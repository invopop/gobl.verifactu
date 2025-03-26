package doc

import (
	"errors"
	"fmt"
	"time"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var correctiveCodes = []cbc.Code{ // Credit or Debit notes
	"R1", "R2", "R3", "R4", "R5",
}

// newInvoice creates a new VeriFactu registration for an invoice.
func newInvoice(inv *bill.Invoice, ts time.Time, r IssuerRole, s *Software) (*RegistroAlta, error) {
	tf, err := getTaxExtKey(inv, verifactu.ExtKeyDocType)
	if err != nil {
		return nil, err
	}

	desc, err := newDescription(inv.Notes)
	if err != nil {
		return nil, err
	}

	dg, err := newDesglose(inv)
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
		TipoFactura:              tf,
		DescripcionOperacion:     desc,
		Desglose:                 dg,
		CuotaTotal:               newTotalTaxes(inv).String(),
		ImporteTotal:             newImporteTotal(inv).String(),
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

	if inv.Tax.Ext[verifactu.ExtKeyDocType].In(correctiveCodes...) {
		k, err := getTaxExtKey(inv, verifactu.ExtKeyCorrectionType)
		if err != nil {
			return nil, err
		}
		reg.TipoRectificativa = k

		list := make([]*FacturaRectificada, len(inv.Preceding))
		taxes := new(tax.Total)
		for i, ref := range inv.Preceding {
			if ref.Tax != nil {
				taxes = taxes.Merge(ref.Tax)
			}
			list[i] = &FacturaRectificada{
				IDFactura: IDFactura{
					IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
					NumSerieFactura:        invoiceNumber(ref.Series, ref.Code),
					FechaExpedicionFactura: ref.IssueDate.Time().Format("02-01-2006"),
				},
			}
		}
		reg.FacturasRectificadas = list
		reg.ImporteRectificacion = newImporteRectificacion(taxes)
	}

	// F3 covers the special use-case of full invoices that replace a
	// previous simplified document. This is the only time the "FacturaSustituida"
	// field is used.
	if reg.TipoFactura == "F3" {
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

func newDescription(notes []*org.Note) (string, error) {
	for _, note := range notes {
		if note.Key == org.NoteKeyGeneral {
			return note.Text, nil
		}
	}
	return "", ErrValidation.WithMessage(fmt.Sprintf("notes: missing note with key '%s'", org.NoteKeyGeneral))
}

func newImporteTotal(inv *bill.Invoice) num.Amount {
	totalWithDiscounts := inv.Totals.Total
	if inv.Totals.Taxes == nil {
		return totalWithDiscounts
	}
	totalTaxes := num.MakeAmount(0, 2)
	for _, category := range inv.Totals.Taxes.Categories {
		if !category.Retained {
			totalTaxes = totalTaxes.Add(category.Amount)
		}
	}
	return totalWithDiscounts.Add(totalTaxes)
}

func newImporteRectificacion(taxes *tax.Total) *ImporteRectificacion {
	zero := currency.EUR.Def().Zero()
	ir := &ImporteRectificacion{
		BaseRectificada:         zero,
		CuotaRectificada:        zero,
		CuotaRecargoRectificado: zero,
	}
	for _, cat := range taxes.Categories {
		if cat.Code == tax.CategoryVAT {
			for _, rate := range cat.Rates {
				ir.BaseRectificada = ir.BaseRectificada.Add(rate.Base)
			}
			ir.CuotaRectificada = ir.CuotaRectificada.Add(cat.Amount)
			if cat.Surcharge != nil {
				ir.CuotaRecargoRectificado = ir.CuotaRecargoRectificado.Add(*cat.Surcharge)
			}
		}
	}
	return ir
}

func newTotalTaxes(inv *bill.Invoice) num.Amount {
	totalTaxes := num.MakeAmount(0, 2)
	if inv.Totals.Taxes == nil {
		return totalTaxes
	}
	for _, category := range inv.Totals.Taxes.Categories {
		if !category.Retained {
			totalTaxes = totalTaxes.Add(category.Amount)
		}
	}
	return totalTaxes
}

func getTaxExtKey(inv *bill.Invoice, k cbc.Key) (string, error) {
	if inv.Tax == nil || inv.Tax.Ext == nil || inv.Tax.Ext[k].String() == "" {
		return "", validation.Errors{
			"tax": validation.Errors{
				"ext": validation.Errors{
					k.String(): errors.New("required"),
				},
			},
		}
	}
	return inv.Tax.Ext[k].String(), nil
}

package verifactu

import (
	"fmt"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
)

// L1 tax codes
const (
	taxCodeVAT   = "01"
	taxCodeIPSI  = "02"
	taxCodeIGIC  = "03"
	taxCodeOther = "05"
)

var taxCategoryCodeMap = map[cbc.Code]string{
	tax.CategoryVAT:    taxCodeVAT,
	es.TaxCategoryIPSI: taxCodeIPSI,
	es.TaxCategoryIGIC: taxCodeIGIC,
}

func newDesglose(inv *bill.Invoice) (*Desglose, error) {
	if inv.Totals == nil || inv.Totals.Taxes == nil {
		return nil, nil
	}

	desglose := &Desglose{}
	for _, c := range inv.Totals.Taxes.Categories {
		if c.Retained {
			continue
		}
		for _, r := range c.Rates {
			detalleDesglose, err := buildDetalleDesglose(c, r)
			if err != nil {
				return nil, err
			}
			desglose.DetalleDesglose = append(desglose.DetalleDesglose, detalleDesglose)
		}
	}

	return desglose, nil
}

// Rules applied to build the breakdown come from:
// https://www.agenciatributaria.es/static_files/AEAT_Desarrolladores/EEDD/IVA/VERI-FACTU/Validaciones_Errores_Veri-Factu.pdf
func buildDetalleDesglose(c *tax.CategoryTotal, r *tax.RateTotal) (*DetalleDesglose, error) {
	detalle := &DetalleDesglose{
		BaseImponibleOImporteNoSujeto: r.Base.String(),
	}

	cat, ok := taxCategoryCodeMap[c.Code]
	if !ok {
		detalle.Impuesto = taxCodeOther
	} else {
		detalle.Impuesto = cat
	}

	if c.Code == tax.CategoryVAT || c.Code == es.TaxCategoryIGIC {
		detalle.ClaveRegimen = r.Ext.Get(verifactu.ExtKeyRegime).String()
	}

	if r.Ext == nil {
		return nil, ErrValidation.WithMessage(fmt.Sprintf("missing tax extensions for rate %s", r.Key))
	}

	if r.Percent == nil && r.Ext.Has(verifactu.ExtKeyExempt) {
		detalle.OperacionExenta = r.Ext[verifactu.ExtKeyExempt].String()
	} else if r.Ext.Has(verifactu.ExtKeyOpClass) {
		detalle.CalificacionOperacion = r.Ext.Get(verifactu.ExtKeyOpClass).String()
		switch detalle.CalificacionOperacion {
		case "S1", "S2":
			// Exempt operations should never show amounts, even if zero
			// S2 implies reverse-charge mechanism, so should be 0.
			detalle.CuotaRepercutida = r.Amount.String()
		}
	}

	if detalle.Impuesto == taxCodeIPSI || detalle.Impuesto == taxCodeOther || detalle.ClaveRegimen == "06" {
		detalle.BaseImponibleACoste = r.Base.String()
	}

	switch detalle.CalificacionOperacion {
	case "S1":
		detalle.TipoImpositivo = r.Percent.StringWithoutSymbol()

		// Surcharges can only happen for regular national transactions.
		if r.Surcharge != nil {
			detalle.TipoRecargoEquivalencia = r.Surcharge.Percent.StringWithoutSymbol()
			detalle.CuotaRecargoEquivalencia = r.Surcharge.Amount.String()
		}
	case "S2":
		// Implies reverse-charge with 0 rate
		detalle.TipoImpositivo = "0"
	}

	if detalle.OperacionExenta == "" && detalle.CalificacionOperacion == "" {
		return nil, ErrValidation.WithMessage(fmt.Sprintf("missing operation classification for rate %s", r.Key))
	}

	return detalle, nil
}

package doc

import (
	"fmt"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
)

var taxCategoryCodeMap = map[cbc.Code]string{
	tax.CategoryVAT:    "01",
	es.TaxCategoryIGIC: "02",
	es.TaxCategoryIPSI: "03",
}

func newDesglose(inv *bill.Invoice) (*Desglose, error) {
	desglose := &Desglose{}

	for _, c := range inv.Totals.Taxes.Categories {
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
		BaseImponibleOImporteNoSujeto: r.Base.Float64(),
		CuotaRepercutida:              r.Amount.Float64(),
	}

	cat, ok := taxCategoryCodeMap[c.Code]
	if !ok {
		detalle.Impuesto = "05"
	} else {
		detalle.Impuesto = cat
	}

	if c.Code == tax.CategoryVAT || c.Code == es.TaxCategoryIGIC {
		detalle.ClaveRegimen = r.Ext.Get(verifactu.ExtKeyRegime).String()
	}

	if r.Ext == nil {
		return nil, fmt.Errorf("missing tax extensions for rate %s", r.Key)
	}

	if r.Percent == nil {
		detalle.OperacionExenta = r.Ext[verifactu.ExtKeyOpClass].String()
	} else if !r.Ext.Has(verifactu.ExtKeyOpClass) {
		detalle.CalificacionOperacion = r.Ext.Get(verifactu.ExtKeyExempt).String()
	}

	if detalle.Impuesto == "02" || detalle.Impuesto == "05" || detalle.ClaveRegimen == "06" {
		detalle.BaseImponibleACoste = r.Base.Float64()
	}

	if r.Percent != nil {
		detalle.TipoImpositivo = r.Percent.Amount().Float64()
	}

	if detalle.OperacionExenta == "" && detalle.CalificacionOperacion == "" {
		return nil, fmt.Errorf("missing operation classification for rate %s", r.Key)
	}

	if r.Key.Has(es.TaxRateEquivalence) {
		detalle.TipoRecargoEquivalencia = r.Surcharge.Percent.Amount().Float64()
		detalle.CuotaRecargoEquivalencia = r.Surcharge.Amount.Float64()
	}

	return detalle, nil
}

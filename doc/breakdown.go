package doc

import (
	"fmt"
	"strings"

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

	// L1: IVA, IGIC, IPSI or Other. Default is IVA.
	if r.Ext != nil && r.Ext[verifactu.ExtKeyTaxCategory] != "" {
		cat, ok := taxCategoryCodeMap[c.Code]
		if !ok {
			detalle.Impuesto = "05"
		} else {
			detalle.Impuesto = cat
		}
		if detalle.Impuesto == "01" {
			// L8A: IVA
			detalle.ClaveRegimen = r.Ext[verifactu.ExtKeyTaxRegime].String()
		}
		if detalle.Impuesto == "02" {
			// L8B: IGIC
			detalle.ClaveRegimen = r.Ext[verifactu.ExtKeyTaxRegime].String()
		}
	} else {
		// L8A: IVA
		detalle.ClaveRegimen = r.Ext[verifactu.ExtKeyTaxRegime].String()
	}

	// Rate zero is what VeriFactu calls "Exempt operation", in difference to GOBL's exempt operation, which in
	// VeriFactu is called "No sujeta".
	if r.Key == tax.RateZero {
		detalle.OperacionExenta = r.Ext[verifactu.ExtKeyTaxClassification].String()
		if detalle.OperacionExenta != "" && !strings.HasPrefix(detalle.OperacionExenta, "E") {
			return nil, fmt.Errorf("invalid exemption code %s - must be E1-E6", detalle.OperacionExenta)
		}
	} else {
		detalle.CalificacionOperacion = r.Ext[verifactu.ExtKeyTaxClassification].String()
		if detalle.CalificacionOperacion == "" {
			return nil, fmt.Errorf("missing operation classification for rate %s", r.Key)
		}
	}

	if isSpecialRegime(c, r) {
		detalle.BaseImponibleACoste = r.Base.Float64()
	}

	if r.Percent != nil {
		detalle.TipoImpositivo = r.Percent.Amount().Float64()
	}

	if detalle.OperacionExenta == "" && detalle.CalificacionOperacion == "" {
		return nil, fmt.Errorf("missing operation classification for rate %s", r.Key)
	}

	if hasEquivalenceSurcharge(c, r) {
		if r.Surcharge == nil {
			return nil, fmt.Errorf("missing surcharge for rate %s", r.Key)
		}
		detalle.TipoRecargoEquivalencia = r.Surcharge.Percent.Amount().Float64()
		detalle.CuotaRecargoEquivalencia = r.Surcharge.Amount.Float64()
	}

	return detalle, nil
}

func isSpecialRegime(c *tax.CategoryTotal, r *tax.RateTotal) bool {
	return r.Ext != nil && (c.Code == es.TaxCategoryIGIC || c.Code == es.TaxCategoryIPSI || r.Ext[verifactu.ExtKeyTaxRegime] == "18")
}

func hasEquivalenceSurcharge(c *tax.CategoryTotal, r *tax.RateTotal) bool {
	return r.Ext != nil && c.Code == tax.CategoryVAT && r.Ext[verifactu.ExtKeyTaxRegime] == "18"
}

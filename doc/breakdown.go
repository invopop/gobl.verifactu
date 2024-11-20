package doc

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
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
			detalleDesglose, err := buildDetalleDesglose(inv, c, r)
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
func buildDetalleDesglose(inv *bill.Invoice, c *tax.CategoryTotal, r *tax.RateTotal) (*DetalleDesglose, error) {
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
		detalle.ClaveRegimen = parseClave(inv, c, r)
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

func parseClave(inv *bill.Invoice, c *tax.CategoryTotal, r *tax.RateTotal) string {
	switch c.Code {
	case tax.CategoryVAT:
		if inv.Customer != nil && partyTaxCountry(inv.Customer) != "ES" {
			return "02"
		}
		if inv.HasTags(es.TagSecondHandGoods) || inv.HasTags(es.TagAntiques) || inv.HasTags(es.TagArt) {
			return "03"
		}
		if inv.HasTags(es.TagTravelAgency) {
			return "05"
		}
		if r.Key == es.TaxRateEquivalence {
			return "18"
		}
		if inv.HasTags(es.TagSimplifiedScheme) {
			return "20"
		}
		return "01"
	case es.TaxCategoryIGIC:
		if inv.Customer != nil && partyTaxCountry(inv.Customer) != "ES" {
			return "02"
		}
		if inv.HasTags(es.TagSecondHandGoods) || inv.HasTags(es.TagAntiques) || inv.HasTags(es.TagArt) {
			return "03"
		}
		if inv.HasTags(es.TagTravelAgency) {
			return "05"
		}
		return "01"
	}
	return ""
}

func partyTaxCountry(party *org.Party) l10n.TaxCountryCode {
	if party != nil && party.TaxID != nil {
		return party.TaxID.Country
	}
	return "ES"
}

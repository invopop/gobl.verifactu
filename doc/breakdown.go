package doc

import (
	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
)

func newDesglose(inv *bill.Invoice) (*Desglose, error) {
	desglose := &Desglose{}

	for _, c := range inv.Totals.Taxes.Categories {
		for _, r := range c.Rates {
			detalleDesglose, err := buildDetalleDesglose(r)
			if err != nil {
				return nil, err
			}
			desglose.DetalleDesglose = append(desglose.DetalleDesglose, detalleDesglose)
		}
	}

	return desglose, nil
}

func buildDetalleDesglose(r *tax.RateTotal) (*DetalleDesglose, error) {
	detalle := &DetalleDesglose{
		BaseImponibleOImporteNoSujeto: r.Base.String(),
		CuotaRepercutida:              r.Amount.String(),
	}

	if r.Ext != nil && r.Ext[verifactu.ExtKeyTaxCategory] != "" {
		detalle.Impuesto = r.Ext[verifactu.ExtKeyTaxCategory].String()
	}

	if r.Key == tax.RateExempt {
		detalle.OperacionExenta = r.Ext[verifactu.ExtKeyExemption].String()
	} else {
		detalle.CalificacionOperacion = r.Ext[verifactu.ExtKeyTaxClassification].String()
	}

	if r.Percent != nil {
		detalle.TipoImpositivo = r.Percent.String()
	}
	return detalle, nil
}

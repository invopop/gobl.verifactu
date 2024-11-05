package doc

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// DesgloseIVA contains the breakdown of VAT amounts
type DesgloseIVA struct {
	DetalleIVA []*DetalleIVA
}

// DetalleIVA details about taxed amounts
type DetalleIVA struct {
	TipoImpositivo    string
	BaseImponible     string
	CuotaImpuesto     string
	TipoRecargoEquiv  string
	CuotaRecargoEquiv string
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

func buildDetalleDesglose(c *tax.CategoryTotal, r *tax.RateTotal) (*DetalleDesglose, error) {
	detalle := &DetalleDesglose{
		BaseImponibleOImporteNoSujeto: r.Base.String(),
		CuotaRepercutida:              r.Amount.String(),
	}

	// MAL - mapear a codigo
	if c.Code != cbc.CodeEmpty {
		detalle.Impuesto = c.Code.String()
	}

	if r.Key == tax.RateExempt {
		detalle.OperacionExenta = "1"
	} else {
		detalle.CalificacionOperacion = "1"
	}

	if r.Percent != nil {
		detalle.TipoImpositivo = r.Percent.String()
	}
	return detalle, nil
}

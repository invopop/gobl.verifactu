package doc_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/gobl.verifactu/test"
	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func defaultBreakdownOptions() *doc.Options {
	return &doc.Options{
		Software:   nil,
		IssuerRole: doc.IssuerRoleSupplier,
		Timestamp:  time.Now(),
	}
}

func TestBreakdownConversion(t *testing.T) {
	t.Run("basic-invoice", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		_ = inv.Calculate()
		d, err := doc.NewInvoice(inv, defaultBreakdownOptions())
		require.NoError(t, err)

		assert.Equal(t, "1800.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "378.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CuotaRepercutida)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].Impuesto)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].ClaveRegimen)
		assert.Equal(t, "S1", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CalificacionOperacion)
	})

	t.Run("exempt-taxes", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Price: num.MakeAmount(100, 0),
				},
				Taxes: tax.Set{
					&tax.Combo{
						Category: "VAT",
						Ext: tax.Extensions{
							verifactu.ExtKeyExempt: "E1",
							verifactu.ExtKeyRegime: "01",
						},
					},
				},
			},
		}
		_ = inv.Calculate()
		d, err := doc.NewInvoice(inv, defaultBreakdownOptions())
		require.NoError(t, err)
		assert.Equal(t, "100.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].Impuesto)
		assert.Equal(t, "E1", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].OperacionExenta)
	})

	t.Run("multiple-tax-rates", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Price: num.MakeAmount(100, 0),
				},
				Taxes: tax.Set{
					&tax.Combo{
						Category: "VAT",
						Rate:     "standard",
						Ext: tax.Extensions{
							verifactu.ExtKeyOpClass: "S1",
							verifactu.ExtKeyRegime:  "01",
						},
					},
				},
			},
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Price: num.MakeAmount(50, 0),
				},
				Taxes: tax.Set{
					&tax.Combo{
						Category: "VAT",
						Rate:     "reduced",
						Ext: tax.Extensions{
							verifactu.ExtKeyOpClass: "S1",
						},
					},
				},
			},
		}
		_ = inv.Calculate()
		d, err := doc.NewInvoice(inv, defaultBreakdownOptions())
		require.NoError(t, err)
		assert.Equal(t, "100.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "21.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CuotaRepercutida)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].Impuesto)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].ClaveRegimen)

		assert.Equal(t, "50.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[1].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "5.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[1].CuotaRepercutida)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[1].Impuesto)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[1].ClaveRegimen)
	})

	t.Run("not-subject-taxes", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Price: num.MakeAmount(100, 0),
				},
				Taxes: tax.Set{
					&tax.Combo{
						Category: "VAT",
						Rate:     "exempt",
						Ext: tax.Extensions{
							verifactu.ExtKeyOpClass: "N1",
							verifactu.ExtKeyRegime:  "01",
						},
					},
				},
			},
		}
		_ = inv.Calculate()
		d, err := doc.NewInvoice(inv, defaultBreakdownOptions())
		require.NoError(t, err)
		assert.Equal(t, "100.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "0.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CuotaRepercutida)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].Impuesto)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].ClaveRegimen)
		assert.Equal(t, "N1", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CalificacionOperacion)
	})

	t.Run("equivalence-surcharge", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Price: num.MakeAmount(100, 0),
				},
				Taxes: tax.Set{
					&tax.Combo{
						Category: "VAT",
						Rate:     "standard+eqs",
						Ext: tax.Extensions{
							verifactu.ExtKeyOpClass: "S1",
							verifactu.ExtKeyRegime:  "01",
						},
					},
				},
			},
		}
		_ = inv.Calculate()
		d, err := doc.NewInvoice(inv, defaultBreakdownOptions())
		require.NoError(t, err)
		assert.Equal(t, "100.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "21.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CuotaRepercutida)
		assert.Equal(t, "5.20", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CuotaRecargoEquivalencia)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].Impuesto)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].ClaveRegimen)
		assert.Equal(t, "S1", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CalificacionOperacion)
	})

	t.Run("ipsi-tax", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		p := num.MakePercentage(10, 2)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Price: num.MakeAmount(100, 0),
				},
				Taxes: tax.Set{
					&tax.Combo{
						Category: es.TaxCategoryIPSI,
						Percent:  &p,
						Ext: tax.Extensions{
							verifactu.ExtKeyOpClass: "S1",
						},
					},
				},
			},
		}
		_ = inv.Calculate()
		d, err := doc.NewInvoice(inv, defaultBreakdownOptions())
		require.NoError(t, err)
		assert.Equal(t, "100.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "10.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CuotaRepercutida)
		assert.Equal(t, "03", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].Impuesto)
		assert.Empty(t, d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].ClaveRegimen)
		assert.Equal(t, "S1", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CalificacionOperacion)
	})

	t.Run("antiques", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Price: num.MakeAmount(1000, 0),
				},
				Taxes: tax.Set{
					&tax.Combo{
						Category: "VAT",
						Rate:     "reduced",
						Ext: tax.Extensions{
							verifactu.ExtKeyOpClass: "S1",
							verifactu.ExtKeyRegime:  "04",
						},
					},
				},
			},
		}
		_ = inv.Calculate()
		d, err := doc.NewInvoice(inv, defaultBreakdownOptions())
		require.NoError(t, err)
		assert.Equal(t, "1000.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "100.00", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CuotaRepercutida)
		assert.Equal(t, "01", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].Impuesto)
		assert.Equal(t, "04", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].ClaveRegimen)
		assert.Equal(t, "S1", d.Body.VeriFactu.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CalificacionOperacion)
	})
}

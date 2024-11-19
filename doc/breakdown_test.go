package doc

import (
	"testing"
	"time"

	"github.com/invopop/gobl.verifactu/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBreakdownConversion(t *testing.T) {
	t.Run("should handle basic invoice breakdown", func(t *testing.T) {
		inv := test.LoadInvoice("./test/data/inv-base.json")
		err := inv.Calculate()
		require.NoError(t, err)

		doc, err := NewDocument(inv, time.Now(), IssuerRoleSupplier, nil)
		require.NoError(t, err)

		assert.Equal(t, 200.00, doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, 42.00, doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CuotaRepercutida)
		assert.Equal(t, "01", doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].Impuesto)
		assert.Equal(t, "01", doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].ClaveRegimen)
	})

	t.Run("should handle exempt taxes", func(t *testing.T) {
		inv := test.LoadInvoice("./test/data/inv-base.json")
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Price: num.MakeAmount(100, 0),
				},
				Taxes: tax.Set{
					&tax.Combo{
						Category: "VAT",
						Rate:     "zero",
					},
				},
			},
		}
		err := inv.Calculate()
		require.NoError(t, err)

		doc, err := NewDocument(inv, time.Now(), IssuerRoleSupplier, nil)
		require.NoError(t, err)

		assert.Equal(t, 100.00, doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "", doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].Impuesto)
	})

	t.Run("should handle multiple tax rates", func(t *testing.T) {
		inv := test.LoadInvoice("./test/data/inv-base.json")
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
					},
				},
			},
		}
		_ = inv.Calculate()

		doc, err := NewDocument(inv, time.Now(), IssuerRoleSupplier, nil)
		require.NoError(t, err)

		assert.Equal(t, 100.00, doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, 21.00, doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].CuotaRepercutida)
		assert.Equal(t, "01", doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[0].ClaveRegimen)

		assert.Equal(t, 50.00, doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[1].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, 10.50, doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[1].CuotaRepercutida)
		assert.Equal(t, "01", doc.RegistroFactura.RegistroAlta.Desglose.DetalleDesglose[1].ClaveRegimen)

	})
}

package doc_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/gobl.verifactu/test"
	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistroAlta(t *testing.T) {
	ts, err := time.Parse(time.RFC3339, "2022-02-01T04:00:00Z")
	require.NoError(t, err)
	opts := &doc.Options{
		Software:   &doc.Software{},
		IssuerRole: doc.IssuerRoleSupplier,
		Timestamp:  ts,
	}

	t.Run("should contain basic document info", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		d, err := doc.NewInvoice(inv, opts)
		require.NoError(t, err)

		reg := d.Body.VeriFactu.RegistroFactura.RegistroAlta
		assert.Equal(t, "1.0", reg.IDVersion)
		assert.Equal(t, "B85905495", reg.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-004", reg.IDFactura.NumSerieFactura)
		assert.Equal(t, "13-11-2024", reg.IDFactura.FechaExpedicionFactura)
		assert.Equal(t, "Invopop S.L.", reg.NombreRazonEmisor)
		assert.Equal(t, "F1", reg.TipoFactura)
		assert.Equal(t, "This is a sample invoice with a standard tax", reg.DescripcionOperacion)
		assert.Equal(t, "378.00", reg.CuotaTotal)
		assert.Equal(t, "2178.00", reg.ImporteTotal)

		require.Len(t, reg.Destinatarios, 1)
		dest := reg.Destinatarios[0].IDDestinatario
		assert.Equal(t, "Sample Consumer", dest.NombreRazon)
		assert.Equal(t, "B63272603", dest.NIF)

		require.Len(t, reg.Desglose.DetalleDesglose, 1)
		desg := reg.Desglose.DetalleDesglose[0]
		assert.Equal(t, "01", desg.Impuesto)
		assert.Equal(t, "01", desg.ClaveRegimen)
		assert.Equal(t, "S1", desg.CalificacionOperacion)
		assert.Equal(t, "21.0", desg.TipoImpositivo)
		assert.Equal(t, "1800.00", desg.BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "378.00", desg.CuotaRepercutida)
	})
	t.Run("should handle simplified invoices", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil

		d, err := doc.NewInvoice(inv, opts)
		require.NoError(t, err)

		assert.Equal(t, "S", d.Body.VeriFactu.RegistroFactura.RegistroAlta.FacturaSinIdentifDestinatarioArt61d)
	})

	t.Run("should handle rectificative invoices", func(t *testing.T) {
		inv := test.LoadInvoice("cred-note-base.json")

		d, err := doc.NewInvoice(inv, opts)
		require.NoError(t, err)

		reg := d.Body.VeriFactu.RegistroFactura.RegistroAlta
		assert.Equal(t, "R1", reg.TipoFactura)
		assert.Equal(t, "I", reg.TipoRectificativa)
		require.Len(t, reg.FacturasRectificadas, 1)

		rectified := reg.FacturasRectificadas[0]
		assert.Equal(t, "B85905495", rectified.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-085", rectified.IDFactura.NumSerieFactura)
		assert.Equal(t, "10-01-2022", rectified.IDFactura.FechaExpedicionFactura)
		assert.Equal(t, "-1620.00", reg.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "-340.20", reg.Desglose.DetalleDesglose[0].CuotaRepercutida)
		assert.Equal(t, "-340.20", reg.CuotaTotal)
		assert.Equal(t, "-1960.20", reg.ImporteTotal)
	})

	t.Run("should handle substitution invoices", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		inv.Preceding = []*org.DocumentRef{
			{
				Series:    "SAMPLE",
				Code:      "002",
				IssueDate: cal.NewDate(2024, 1, 15),
			},
		}
		inv.Tax.Ext[verifactu.ExtKeyDocType] = "F3"

		d, err := doc.NewInvoice(inv, opts)
		require.NoError(t, err)

		reg := d.Body.VeriFactu.RegistroFactura.RegistroAlta
		require.Len(t, reg.FacturasSustituidas, 1)

		substituted := reg.FacturasSustituidas[0]
		assert.Equal(t, "B85905495", substituted.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-002", substituted.IDFactura.NumSerieFactura)
		assert.Equal(t, "15-01-2024", substituted.IDFactura.FechaExpedicionFactura)
	})
}

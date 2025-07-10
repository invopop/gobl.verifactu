package verifactu_test

import (
	"testing"
	"time"

	verifactu "github.com/invopop/gobl.verifactu"
	"github.com/invopop/gobl.verifactu/test"
	addon "github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistroAlta(t *testing.T) {
	ts, err := time.Parse(time.RFC3339, "2022-02-01T04:00:00Z")
	require.NoError(t, err)
	vc, err := verifactu.New(
		nil, // no software
		verifactu.WithCurrentTime(ts),
	)
	require.NoError(t, err)

	t.Run("should contain basic document info", func(t *testing.T) {
		env := test.LoadEnvelope("inv-base.json")
		ra, err := vc.RegisterInvoice(env, nil)
		require.NoError(t, err)

		assert.Equal(t, "1.0", ra.IDVersion)
		assert.Equal(t, "B85905495", ra.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-004", ra.IDFactura.NumSerieFactura)
		assert.Equal(t, "13-11-2024", ra.IDFactura.FechaExpedicionFactura)
		assert.Equal(t, "Invopop S.L.", ra.NombreRazonEmisor)
		assert.Equal(t, "F1", ra.TipoFactura)
		assert.Equal(t, "This is a sample invoice with a standard tax", ra.DescripcionOperacion)
		assert.Equal(t, "378.00", ra.CuotaTotal)
		assert.Equal(t, "2178.00", ra.ImporteTotal)

		require.Len(t, ra.Destinatarios, 1)
		dest := ra.Destinatarios[0].IDDestinatario
		assert.Equal(t, "Sample Consumer", dest.NombreRazon)
		assert.Equal(t, "B63272603", dest.NIF)

		require.Len(t, ra.Desglose.DetalleDesglose, 1)
		desg := ra.Desglose.DetalleDesglose[0]
		assert.Equal(t, "01", desg.Impuesto)
		assert.Equal(t, "01", desg.ClaveRegimen)
		assert.Equal(t, "S1", desg.CalificacionOperacion)
		assert.Equal(t, "21.0", desg.TipoImpositivo)
		assert.Equal(t, "1800.00", desg.BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "378.00", desg.CuotaRepercutida)
	})
	t.Run("should handle simplified invoices", func(t *testing.T) {
		env, inv := test.LoadInvoice("inv-base.json")
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())

		ra, err := vc.RegisterInvoice(env, nil)
		require.NoError(t, err)

		assert.Equal(t, "S", ra.FacturaSinIdentifDestinatarioArt61d)
	})

	t.Run("should handle rectificative invoices", func(t *testing.T) {
		env := test.LoadEnvelope("cred-note-base.json")

		ra, err := vc.RegisterInvoice(env, nil)
		require.NoError(t, err)

		assert.Equal(t, "R1", ra.TipoFactura)
		assert.Equal(t, "I", ra.TipoRectificativa)
		require.Len(t, ra.FacturasRectificadas, 1)

		rectified := ra.FacturasRectificadas[0]
		assert.Equal(t, "B85905495", rectified.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-085", rectified.IDFactura.NumSerieFactura)
		assert.Equal(t, "10-01-2022", rectified.IDFactura.FechaExpedicionFactura)
		assert.Equal(t, "-1620.00", ra.Desglose.DetalleDesglose[0].BaseImponibleOImporteNoSujeto)
		assert.Equal(t, "-340.20", ra.Desglose.DetalleDesglose[0].CuotaRepercutida)
		assert.Equal(t, "-340.20", ra.CuotaTotal)
		assert.Equal(t, "-1960.20", ra.ImporteTotal)
	})

	t.Run("should handle substitution invoices", func(t *testing.T) {
		env, inv := test.LoadInvoice("inv-base.json")
		inv.Preceding = []*org.DocumentRef{
			{
				Series:    "SAMPLE",
				Code:      "002",
				IssueDate: cal.NewDate(2024, 1, 15),
			},
		}
		inv.Tax.Ext[addon.ExtKeyDocType] = "F3"

		ra, err := vc.RegisterInvoice(env, nil)
		require.NoError(t, err)

		require.Len(t, ra.FacturasSustituidas, 1)
		substituted := ra.FacturasSustituidas[0]
		assert.Equal(t, "B85905495", substituted.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-002", substituted.IDFactura.NumSerieFactura)
		assert.Equal(t, "15-01-2024", substituted.IDFactura.FechaExpedicionFactura)
	})

	t.Run("should handle an empty note", func(t *testing.T) {
		env, inv := test.LoadInvoice("inv-base.json")
		inv.Notes = nil
		ra, err := vc.RegisterInvoice(env, nil)
		require.NoError(t, err)

		assert.Equal(t, "1.0", ra.IDVersion)
		assert.Equal(t, "B85905495", ra.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-004", ra.IDFactura.NumSerieFactura)
		assert.Equal(t, "13-11-2024", ra.IDFactura.FechaExpedicionFactura)
		assert.Equal(t, "Invopop S.L.", ra.NombreRazonEmisor)
		assert.Equal(t, "Development services.", ra.DescripcionOperacion)
	})

	t.Run("should handle an empty note with multiple items", func(t *testing.T) {
		env, inv := test.LoadInvoice("inv-base.json")
		inv.Notes = nil
		inv.Lines[0].Item.Name = "a"
		for range int(100) {
			inv.Lines = append(inv.Lines, inv.Lines[0])
		}
		ra, err := vc.RegisterInvoice(env, nil)
		require.NoError(t, err)

		assert.Equal(t, "1.0", ra.IDVersion)
		assert.Equal(t, "B85905495", ra.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-004", ra.IDFactura.NumSerieFactura)
		assert.Equal(t, "13-11-2024", ra.IDFactura.FechaExpedicionFactura)
		assert.Equal(t, "Invopop S.L.", ra.NombreRazonEmisor)
		assert.True(t, len(ra.DescripcionOperacion) <= 500)
	})
}

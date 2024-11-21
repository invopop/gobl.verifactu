package doc

import (
	"testing"
	"time"

	"github.com/invopop/gobl.verifactu/test"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistroAlta(t *testing.T) {
	ts, err := time.Parse(time.RFC3339, "2022-02-01T04:00:00Z")
	require.NoError(t, err)
	role := IssuerRoleSupplier
	sw := &Software{}

	t.Run("should contain basic document info", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		doc, err := NewDocument(inv, ts, role, sw)
		require.NoError(t, err)

		reg := doc.RegistroFactura.RegistroAlta
		assert.Equal(t, "1.0", reg.IDVersion)
		assert.Equal(t, "B85905495", reg.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-003", reg.IDFactura.NumSerieFactura)
		assert.Equal(t, "13-11-2024", reg.IDFactura.FechaExpedicionFactura)
		assert.Equal(t, "Invopop S.L.", reg.NombreRazonEmisor)
		assert.Equal(t, "F1", reg.TipoFactura)
		assert.Equal(t, "This is a sample invoice", reg.DescripcionOperacion)
		assert.Equal(t, float64(378), reg.CuotaTotal)
		assert.Equal(t, float64(2178), reg.ImporteTotal)

		require.Len(t, reg.Destinatarios, 1)
		dest := reg.Destinatarios[0].IDDestinatario
		assert.Equal(t, "Sample Consumer", dest.NombreRazon)
		assert.Equal(t, "B63272603", dest.NIF)

		require.Len(t, reg.Desglose.DetalleDesglose, 1)
		desg := reg.Desglose.DetalleDesglose[0]
		assert.Equal(t, "01", desg.Impuesto)
		assert.Equal(t, "01", desg.ClaveRegimen)
		assert.Equal(t, "S1", desg.CalificacionOperacion)
		assert.Equal(t, float64(21), desg.TipoImpositivo)
		assert.Equal(t, float64(1800), desg.BaseImponibleOImporteNoSujeto)
		assert.Equal(t, float64(378), desg.CuotaRepercutida)
	})
	t.Run("should handle simplified invoices", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil

		doc, err := NewDocument(inv, ts, role, sw)
		require.NoError(t, err)

		assert.Equal(t, "S", doc.RegistroFactura.RegistroAlta.FacturaSinIdentifDestinatarioArt61d)
	})
}

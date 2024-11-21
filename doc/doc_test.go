package doc

import (
	"testing"
	"time"

	"github.com/invopop/gobl.verifactu/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceConversion(t *testing.T) {
	ts, err := time.Parse(time.RFC3339, "2022-02-01T04:00:00Z")
	require.NoError(t, err)
	role := IssuerRoleSupplier
	sw := &Software{}

	t.Run("should contain basic document info", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		doc, err := NewDocument(inv, ts, role, sw)

		require.NoError(t, err)
		assert.Equal(t, "Invopop S.L.", doc.Cabecera.Obligado.NombreRazon)
		assert.Equal(t, "B85905495", doc.Cabecera.Obligado.NIF)
		assert.Equal(t, "1.0", doc.RegistroFactura.RegistroAlta.IDVersion)
		assert.Equal(t, "B85905495", doc.RegistroFactura.RegistroAlta.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-003", doc.RegistroFactura.RegistroAlta.IDFactura.NumSerieFactura)
		assert.Equal(t, "13-11-2024", doc.RegistroFactura.RegistroAlta.IDFactura.FechaExpedicionFactura)
	})
}

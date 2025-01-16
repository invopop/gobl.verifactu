package doc_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/gobl.verifactu/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceConversion(t *testing.T) {
	ts, err := time.Parse(time.RFC3339, "2022-02-01T04:00:00Z")
	require.NoError(t, err)
	role := doc.IssuerRoleSupplier
	sw := &doc.Software{}

	t.Run("should contain basic document info", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		opts := &doc.Options{
			Software:   sw,
			IssuerRole: role,
			Timestamp:  ts,
		}
		doc, err := doc.NewInvoice(inv, opts)

		require.NoError(t, err)
		assert.Equal(t, "Invopop S.L.", doc.Body.VeriFactu.Cabecera.Obligado.NombreRazon)
		assert.Equal(t, "B85905495", doc.Body.VeriFactu.Cabecera.Obligado.NIF)
		assert.Equal(t, "1.0", doc.Body.VeriFactu.RegistroFactura.RegistroAlta.IDVersion)
		assert.Equal(t, "B85905495", doc.Body.VeriFactu.RegistroFactura.RegistroAlta.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-004", doc.Body.VeriFactu.RegistroFactura.RegistroAlta.IDFactura.NumSerieFactura)
		assert.Equal(t, "13-11-2024", doc.Body.VeriFactu.RegistroFactura.RegistroAlta.IDFactura.FechaExpedicionFactura)
	})
}

func TestInvoiceConversionWithRep(t *testing.T) {
	ts, err := time.Parse(time.RFC3339, "2022-02-01T04:00:00Z")
	require.NoError(t, err)
	opts := &doc.Options{
		Software:   &doc.Software{},
		IssuerRole: doc.IssuerRoleThirdParty,
		Timestamp:  ts,
		Representative: &doc.Obligado{
			NombreRazon: "Sample Rep",
			NIF:         "B63272603",
		},
	}

	t.Run("should contain basic document info", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")

		doc, err := doc.NewInvoice(inv, opts)

		require.NoError(t, err)
		assert.Equal(t, "Invopop S.L.", doc.Body.VeriFactu.Cabecera.Obligado.NombreRazon)
		assert.Equal(t, "B85905495", doc.Body.VeriFactu.Cabecera.Obligado.NIF)
		assert.Equal(t, "Sample Rep", doc.Body.VeriFactu.Cabecera.Representante.NombreRazon)
		assert.Equal(t, "B63272603", doc.Body.VeriFactu.Cabecera.Representante.NIF)
		assert.Equal(t, "1.0", doc.Body.VeriFactu.RegistroFactura.RegistroAlta.IDVersion)
		assert.Equal(t, "B85905495", doc.Body.VeriFactu.RegistroFactura.RegistroAlta.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-004", doc.Body.VeriFactu.RegistroFactura.RegistroAlta.IDFactura.NumSerieFactura)
		assert.Equal(t, "13-11-2024", doc.Body.VeriFactu.RegistroFactura.RegistroAlta.IDFactura.FechaExpedicionFactura)
	})
}

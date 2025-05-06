package verifactu_test

import (
	"testing"
	"time"

	verifactu "github.com/invopop/gobl.verifactu"
	"github.com/invopop/gobl.verifactu/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceConversion(t *testing.T) {
	ts, err := time.Parse(time.RFC3339, "2022-02-01T04:00:00Z")
	require.NoError(t, err)
	vc, err := verifactu.New(
		new(verifactu.Software),
		verifactu.WithSupplierIssuer(),
		verifactu.WithCurrentTime(ts),
	)
	require.NoError(t, err)

	t.Run("should contain basic document info", func(t *testing.T) {
		env := test.LoadEnvelope("inv-base.json")
		irq, err := vc.NewEnvelopeInvoiceRequest(env, nil)
		require.NoError(t, err)

		head := irq.Header
		req := irq.Lines[0].Registration
		assert.Equal(t, "Invopop S.L.", head.Obligado.NombreRazon)
		assert.Equal(t, "B85905495", head.Obligado.NIF)
		assert.Equal(t, "1.0", req.IDVersion)
		assert.Equal(t, "B85905495", req.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-004", req.IDFactura.NumSerieFactura)
		assert.Equal(t, "13-11-2024", req.IDFactura.FechaExpedicionFactura)
	})
}

func TestInvoiceConversionWithRep(t *testing.T) {
	ts, err := time.Parse(time.RFC3339, "2022-02-01T04:00:00Z")
	require.NoError(t, err)

	vc, err := verifactu.New(
		new(verifactu.Software),
		verifactu.WithThirdPartyIssuer(),
		verifactu.WithCurrentTime(ts),
		verifactu.WithRepresentative("Sample Rep", "B63272603"),
	)
	require.NoError(t, err)

	t.Run("should contain basic document info", func(t *testing.T) {
		env := test.LoadEnvelope("inv-base.json")
		irq, err := vc.NewEnvelopeInvoiceRequest(env, nil)
		require.NoError(t, err)

		head := irq.Header
		req := irq.Lines[0].Registration
		assert.Equal(t, "Invopop S.L.", head.Obligado.NombreRazon)
		assert.Equal(t, "B85905495", head.Obligado.NIF)
		assert.Equal(t, "Sample Rep", head.Representante.NombreRazon)
		assert.Equal(t, "B63272603", head.Representante.NIF)
		assert.Equal(t, "1.0", req.IDVersion)
		assert.Equal(t, "B85905495", req.IDFactura.IDEmisorFactura)
		assert.Equal(t, "SAMPLE-004", req.IDFactura.NumSerieFactura)
		assert.Equal(t, "13-11-2024", req.IDFactura.FechaExpedicionFactura)
	})
}

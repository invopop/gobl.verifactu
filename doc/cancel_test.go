package doc_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/gobl.verifactu/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistroAnulacion(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		inv := test.LoadInvoice("cred-note-base.json")
		opts := &doc.Options{
			Software:   nil,
			IssuerRole: doc.IssuerRoleSupplier,
			Timestamp:  time.Now(),
		}
		d, err := doc.NewCancel(inv, opts)
		require.NoError(t, err)

		ra := d.Body.VeriFactu.RegistroFactura.RegistroAnulacion
		assert.Equal(t, "B85905495", ra.IDFactura.IDEmisorFactura)
		assert.Equal(t, "FR-012", ra.IDFactura.NumSerieFactura)
		assert.Equal(t, "01-02-2022", ra.IDFactura.FechaExpedicionFactura)
		assert.Equal(t, "01", ra.TipoHuella)
	})
}

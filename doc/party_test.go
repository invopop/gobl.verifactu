package doc_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/gobl.verifactu/test"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewParty(t *testing.T) {
	ts, err := time.Parse(time.RFC3339, "2022-02-01T04:00:00Z")
	require.NoError(t, err)
	opts := &doc.Options{
		Software:   &doc.Software{},
		IssuerRole: doc.IssuerRoleSupplier,
		Timestamp:  ts,
	}

	t.Run("with tax ID", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		d, err := doc.NewInvoice(inv, opts)
		require.NoError(t, err)

		p := d.Body.VeriFactu.RegistroFactura.RegistroAlta.Destinatarios[0].IDDestinatario
		assert.Equal(t, "Sample Consumer", p.NombreRazon)
		assert.Equal(t, "B63272603", p.NIF)
		assert.Nil(t, p.IDOtro)
	})

	t.Run("with passport", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		inv.Customer = &org.Party{
			Name: "Mr. Pass Port",
			Identities: []*org.Identity{
				{
					Key:  org.IdentityKeyPassport,
					Code: "12345",
				},
			},
		}

		d, err := doc.NewInvoice(inv, opts)
		require.NoError(t, err)

		p := d.Body.VeriFactu.RegistroFactura.RegistroAlta.Destinatarios[0].IDDestinatario
		assert.Equal(t, "Mr. Pass Port", p.NombreRazon)
		assert.Empty(t, p.NIF)
		assert.NotNil(t, p.IDOtro)
		assert.Equal(t, "03", p.IDOtro.IDType)
		assert.Equal(t, "12345", p.IDOtro.ID)
	})

	t.Run("with foreign identity", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		inv.Customer = &org.Party{
			Name: "Foreign Company",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "111111125",
			},
		}

		d, err := doc.NewInvoice(inv, opts)
		require.NoError(t, err)

		p := d.Body.VeriFactu.RegistroFactura.RegistroAlta.Destinatarios[0].IDDestinatario

		assert.Equal(t, "Foreign Company", p.NombreRazon)
		assert.Empty(t, p.NIF)
		assert.NotNil(t, p.IDOtro)
		assert.Equal(t, "04", p.IDOtro.IDType)
		assert.Equal(t, "111111125", p.IDOtro.ID)
	})

	t.Run("with no identifiers", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		inv.Customer = &org.Party{
			Name: "Simple Company",
		}

		_, err := doc.NewInvoice(inv, opts)
		require.Error(t, err)
	})
}

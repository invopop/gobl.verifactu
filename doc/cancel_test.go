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

		d, err := doc.NewDocument(inv, time.Now(), doc.IssuerRoleSupplier, nil, true)
		require.NoError(t, err)

		reg := d.RegistroFactura.RegistroAnulacion

		assert.Equal(t, "E", reg.GeneradoPor)
		assert.NotNil(t, reg.Generador)
		assert.Equal(t, "Provide One S.L.", reg.Generador.NombreRazon)
		assert.Equal(t, "B98602642", reg.Generador.NIF)

		data, err := d.BytesIndent()
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})

	t.Run("customer issuer", func(t *testing.T) {
		inv := test.LoadInvoice("cred-note-base.json")

		d, err := doc.NewDocument(inv, time.Now(), doc.IssuerRoleCustomer, nil, true)
		require.NoError(t, err)

		reg := d.RegistroFactura.RegistroAnulacion

		assert.Equal(t, "D", reg.GeneradoPor)
		assert.NotNil(t, reg.Generador)
		assert.Equal(t, "Sample Customer", reg.Generador.NombreRazon)
		assert.Equal(t, "54387763P", reg.Generador.NIF)

		data, err := d.BytesIndent()
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})

	t.Run("third party issuer", func(t *testing.T) {
		inv := test.LoadInvoice("cred-note-base.json")

		d, err := doc.NewDocument(inv, time.Now(), doc.IssuerRoleThirdParty, nil, true)
		require.NoError(t, err)

		reg := d.RegistroFactura.RegistroAnulacion

		assert.Equal(t, "T", reg.GeneradoPor)
		assert.NotNil(t, reg.Generador)
		assert.Equal(t, "Provide One S.L.", reg.Generador.NombreRazon)
		assert.Equal(t, "B98602642", reg.Generador.NIF)

		data, err := d.BytesIndent()
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})

}

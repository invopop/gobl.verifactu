package doc

import (
	"testing"
	"time"

	"github.com/invopop/gobl.verifactu/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistroAnulacion(t *testing.T) {

	t.Run("basic", func(t *testing.T) {
		inv := test.LoadInvoice("cred-note-base.json")

		doc, err := NewDocument(inv, time.Now(), IssuerRoleSupplier, nil)
		require.NoError(t, err)

		reg, err := NewRegistroAnulacion(inv, time.Now(), IssuerRoleSupplier, nil)
		require.NoError(t, err)

		assert.Equal(t, "E", reg.GeneradoPor)
		assert.NotNil(t, reg.Generador)
		assert.Equal(t, "Provide One S.L.", reg.Generador.NombreRazon)
		assert.Equal(t, "B98602642", reg.Generador.NIF)

		data, err := doc.BytesIndent()
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})

	t.Run("customer issuer", func(t *testing.T) {
		inv := test.LoadInvoice("cred-note-base.json")

		doc, err := NewDocument(inv, time.Now(), IssuerRoleCustomer, nil)
		require.NoError(t, err)

		reg, err := NewRegistroAnulacion(inv, time.Now(), IssuerRoleCustomer, nil)
		require.NoError(t, err)

		assert.Equal(t, "D", reg.GeneradoPor)
		assert.NotNil(t, reg.Generador)
		assert.Equal(t, "Sample Customer", reg.Generador.NombreRazon)
		assert.Equal(t, "54387763P", reg.Generador.NIF)

		data, err := doc.BytesIndent()
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})

	t.Run("third party issuer", func(t *testing.T) {
		inv := test.LoadInvoice("cred-note-base.json")

		doc, err := NewDocument(inv, time.Now(), IssuerRoleThirdParty, nil)
		require.NoError(t, err)

		reg, err := NewRegistroAnulacion(inv, time.Now(), IssuerRoleThirdParty, nil)
		require.NoError(t, err)

		assert.Equal(t, "T", reg.GeneradoPor)
		assert.NotNil(t, reg.Generador)
		assert.Equal(t, "Provide One S.L.", reg.Generador.NombreRazon)
		assert.Equal(t, "B98602642", reg.Generador.NIF)

		data, err := doc.BytesIndent()
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})

}

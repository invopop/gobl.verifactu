package doc

import (
	"testing"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewParty(t *testing.T) {
	t.Run("with tax ID", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Company",
			TaxID: &tax.Identity{
				Code: "B12345678",
			},
		}

		result, err := newParty(party)
		require.NoError(t, err)

		assert.Equal(t, "Test Company", result.NombreRazon)
		assert.Equal(t, "B12345678", result.NIF)
		assert.Nil(t, result.IDOtro)
	})

	t.Run("with identity", func(t *testing.T) {
		party := &org.Party{
			Name: "Foreign Company",
			Identities: []*org.Identity{
				{
					Code: "12345",
					Ext: map[cbc.Key]tax.ExtValue{
						verifactu.ExtKeyIdentity: "02",
					},
				},
			},
		}

		result, err := newParty(party)
		require.NoError(t, err)

		assert.Equal(t, "Foreign Company", result.NombreRazon)
		assert.Empty(t, result.NIF)
		assert.NotNil(t, result.IDOtro)
		assert.Equal(t, "02", result.IDOtro.IDType)
		assert.Equal(t, "12345", result.IDOtro.ID)
	})

	t.Run("with no identifiers", func(t *testing.T) {
		party := &org.Party{
			Name: "Simple Company",
		}

		result, err := newParty(party)
		require.NoError(t, err)

		assert.Equal(t, "Simple Company", result.NombreRazon)
		assert.Empty(t, result.NIF)
		assert.Nil(t, result.IDOtro)
	})
}

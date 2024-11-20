package doc

import (
	"testing"

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
				Country: "ES",
				Code:    "B12345678",
			},
		}

		result, err := newParty(party)
		require.NoError(t, err)

		assert.Equal(t, "Test Company", result.NombreRazon)
		assert.Equal(t, "B12345678", result.NIF)
		assert.Nil(t, result.IDOtro)
	})

	t.Run("with passport", func(t *testing.T) {
		party := &org.Party{
			Name: "Mr. Pass Port",
			Identities: []*org.Identity{
				{
					Key:  org.IdentityKeyPassport,
					Code: "12345",
				},
			},
		}

		result, err := newParty(party)
		require.NoError(t, err)

		assert.Equal(t, "Mr. Pass Port", result.NombreRazon)
		assert.Empty(t, result.NIF)
		assert.NotNil(t, result.IDOtro)
		assert.Equal(t, "03", result.IDOtro.IDType)
		assert.Equal(t, "12345", result.IDOtro.ID)
	})

	t.Run("with foreign identity", func(t *testing.T) {
		party := &org.Party{
			Name: "Foreign Company",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "111111125",
			},
		}

		result, err := newParty(party)
		require.NoError(t, err)

		assert.Equal(t, "Foreign Company", result.NombreRazon)
		assert.Empty(t, result.NIF)
		assert.NotNil(t, result.IDOtro)
		assert.Equal(t, "04", result.IDOtro.IDType)
		assert.Equal(t, "111111125", result.IDOtro.ID)
	})

	t.Run("with no identifiers", func(t *testing.T) {
		party := &org.Party{
			Name: "Simple Company",
		}

		_, err := newParty(party)
		require.Error(t, err)
	})
}

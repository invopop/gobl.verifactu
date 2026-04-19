package test

import (
	"testing"

	"github.com/invopop/xmldsig"
	"github.com/stretchr/testify/require"
)

// Certificate loads the persisted test certificate from test/certs/test.p12.
func Certificate(t *testing.T) *xmldsig.Certificate {
	t.Helper()

	cert, err := xmldsig.LoadCertificate(Path("test", "certs", "test.p12"), "test")
	require.NoError(t, err)

	return cert
}

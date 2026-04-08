package verifactu_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	verifactu "github.com/invopop/gobl.verifactu"
	"github.com/invopop/gobl.verifactu/test"

	"github.com/invopop/xmldsig"
	"github.com/lestrrat-go/libxml2/xsd"
	"github.com/stretchr/testify/require"
)

const (
	msgMissingOutFile    = "output file %s missing, run tests with `--update` flag to create"
	msgUnmatchingOutFile = "output file %s does not match, run tests with `--update` flag to update"
)

func TestXMLGeneration(t *testing.T) {
	examples, err := lookupExamples()
	require.NoError(t, err)

	c := loadClient(t)

	var schema *xsd.Schema
	if *test.UpdateOut {
		schema, err = test.LoadSchema("main.xsd")
		require.NoError(t, err)
	}

	for _, example := range examples {
		name := fmt.Sprintf("should convert %s example file successfully", example)

		t.Run(name, func(t *testing.T) {
			env := test.LoadEnvelope(example)

			// Example Data to Test the Fingerprint.
			prev := &verifactu.ChainData{
				IDIssuer:    "B12345678",
				NumSeries:   "SAMPLE-001",
				IssueDate:   "26-11-2024",
				Fingerprint: "0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF",
			}

			ir, err := c.NewEnvelopeInvoiceRequest(env, prev)
			require.NoError(t, err)

			outPath := test.Path("test", "data", "out",
				strings.TrimSuffix(example, ".json")+".xml",
			)

			data, err := ir.Envelop().BytesIndent()
			require.NoError(t, err)

			if *test.UpdateOut {
				errs := test.ValidateXML(schema, data)
				for _, e := range errs {
					require.NoError(t, e)
				}

				err = os.WriteFile(outPath, data, 0644)
				require.NoError(t, err)
				return
			}

			expected, err := os.ReadFile(outPath)

			require.False(t, os.IsNotExist(err), msgMissingOutFile, filepath.Base(outPath))
			require.NoError(t, err)
			require.Equal(t, string(expected), string(data), msgUnmatchingOutFile, filepath.Base(outPath))
		})
	}
}

func loadClient(t *testing.T) *verifactu.Client {
	t.Helper()

	ts, err := time.Parse(time.RFC3339, "2024-11-26T04:00:00Z")
	require.NoError(t, err)

	cert := test.Certificate(t)

	c, err := verifactu.New(verifactu.Software{
		NombreRazon:                 "My Software",
		NIF:                         "12345678A",
		NombreSistemaInformatico:    "My Software",
		IdSistemaInformatico:        "A1",
		Version:                     "1.0",
		NumeroInstalacion:           "12345678A",
		TipoUsoPosibleSoloVerifactu: "S",
		TipoUsoPosibleMultiOT:       "S",
		IndicadorMultiplesOT:        "N",
	},
		verifactu.WithCurrentTime(ts),
		verifactu.WithCertificate(cert),
		verifactu.WithSigning(
			xmldsig.WithDocID("test-doc-id"),
			xmldsig.WithCurrentTime(func() time.Time { return ts }),
		),
	)
	require.NoError(t, err)

	return c
}

func lookupExamples() ([]string, error) {
	examples, err := filepath.Glob(test.Path("test", "data", "*.json"))
	if err != nil {
		return nil, err
	}

	for i, example := range examples {
		examples[i] = filepath.Base(example)
	}

	return examples, nil
}

package verifactu_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	verifactu "github.com/invopop/gobl.verifactu"
	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/gobl.verifactu/test"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	msgMissingOutFile    = "output file %s missing, run tests with `--update` flag to create"
	msgUnmatchingOutFile = "output file %s does not match, run tests with `--update` flag to update"
)

func TestXMLGeneration(t *testing.T) {
	schema, err := loadSchema()
	require.NoError(t, err)

	examples, err := lookupExamples()
	require.NoError(t, err)

	tbai, err := loadClient()
	require.NoError(t, err)

	for _, example := range examples {
		name := fmt.Sprintf("should convert %s example file successfully", example)

		t.Run(name, func(t *testing.T) {
			data, err := convertExample(tbai, example)
			require.NoError(t, err)

			outPath := test.Path("test", "data", "out",
				strings.TrimSuffix(example, ".json")+".xml",
			)

			if *test.UpdateOut {
				errs := validateDoc(schema, data)
				for _, e := range errs {
					assert.NoError(t, e)
				}
				if len(errs) > 0 {
					assert.Fail(t, "Invalid XML:\n"+string(data))
					return
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

func loadSchema() (*xsd.Schema, error) {
	schemaPath := test.Path("test", "schema", "SuministroLR.xsd")
	schema, err := xsd.ParseFromFile(schemaPath)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func loadClient() (*verifactu.Client, error) {
	ts, err := time.Parse(time.RFC3339, "2024-11-26T04:00:00Z")
	if err != nil {
		return nil, err
	}

	return verifactu.New(&doc.Software{
		NombreRazon:              "My Software",
		NIF:                      "12345678A",
		NombreSistemaInformatico: "My Software",
		IdSistemaInformatico:     "12345678A",
		Version:                  "1.0",
		NumeroInstalacion:        "12345678A",
	},
		verifactu.WithCurrentTime(ts),
		verifactu.WithThirdPartyIssuer(),
	)
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

func convertExample(c *verifactu.Client, example string) ([]byte, error) {
	env := test.LoadEnvelope(example)

	td, err := c.Convert(env)
	if err != nil {
		return nil, err
	}

	err = c.Fingerprint(td, &doc.ChainData{})
	if err != nil {
		return nil, err
	}

	return td.BytesIndent()
}

func validateDoc(schema *xsd.Schema, doc []byte) []error {
	xmlDoc, err := libxml2.ParseString(string(doc))
	if err != nil {
		return []error{err}
	}

	err = schema.Validate(xmlDoc)
	if err != nil {
		return err.(xsd.SchemaValidationError).Errors()
	}

	return nil
}

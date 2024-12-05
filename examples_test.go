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

	c, err := loadClient()
	require.NoError(t, err)

	for _, example := range examples {
		name := fmt.Sprintf("should convert %s example file successfully", example)

		t.Run(name, func(t *testing.T) {
			env := test.LoadEnvelope(example)
			td, err := c.Convert(env)
			require.NoError(t, err)

			// Example Data to Test the Fingerprint.
			prev := &doc.ChainData{
				IDEmisorFactura:        "B12345678",
				NumSerieFactura:        "SAMPLE-001",
				FechaExpedicionFactura: "26-11-2024",
				Huella:                 "0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF",
			}

			err = c.Fingerprint(td, prev)
			require.NoError(t, err)

			outPath := test.Path("test", "data", "out",
				strings.TrimSuffix(example, ".json")+".xml",
			)

			valData, err := td.Bytes()
			require.NoError(t, err)

			valData, err = addNamespaces(valData)
			require.NoError(t, err)

			errs := validateDoc(schema, valData)
			for _, e := range errs {
				assert.NoError(t, e)
			}
			if len(errs) > 0 {
				assert.Fail(t, "Invalid XML:\n"+string(valData))
				return
			}

			if *test.UpdateOut {
				data, err := td.Envelop()
				require.NoError(t, err)

				err = os.WriteFile(outPath, data, 0644)
				require.NoError(t, err)

				return
			}

			expected, err := os.ReadFile(outPath)

			require.False(t, os.IsNotExist(err), msgMissingOutFile, filepath.Base(outPath))
			require.NoError(t, err)
			outData, err := td.Envelop()
			require.NoError(t, err)
			require.Equal(t, string(expected), string(outData), msgUnmatchingOutFile, filepath.Base(outPath))
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

// Helper function to inject namespaces into XML without using Envelop()
// Just for xsd validation purposes
func addNamespaces(data []byte) ([]byte, error) {
	xmlString := string(data)
	xmlNamespaces := ` xmlns:sum="https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroLR.xsd" xmlns:sum1="https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroInformacion.xsd"`
	if !strings.Contains(xmlString, "<sum:RegFactuSistemaFacturacion>") {
		return nil, fmt.Errorf("could not find RegFactuSistemaFacturacion tag in XML")
	}
	xmlString = strings.Replace(xmlString, "<sum:RegFactuSistemaFacturacion>", "<sum:RegFactuSistemaFacturacion"+xmlNamespaces+">", 1)
	finalXMLBytes := []byte(xmlString)
	return finalXMLBytes, nil
}

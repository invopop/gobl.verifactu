package verifactu_test

import (
	"reflect"
	"testing"
	"time"

	vf "github.com/invopop/gobl.verifactu"
	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/gobl.verifactu/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDocument(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *doc.Envelope
		wantErr bool
	}{
		{
			name: "valid document",
			data: []byte(`<?xml version="1.0" encoding="UTF-8"?>
			<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:sum="https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroLR.xsd" xmlns:sum1="https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroInformacion.xsd">
			<soapenv:Body>
				<sum:RegFactuSistemaFacturacion>
				<sum:Cabecera>
					<sum1:ObligadoEmision>
					<sum1:NombreRazon>Test Company</sum1:NombreRazon>
					<sum1:NIF>B12345678</sum1:NIF>
					</sum1:ObligadoEmision>
				</sum:Cabecera>
				</sum:RegFactuSistemaFacturacion>
			</soapenv:Body>
			</soapenv:Envelope>`),
			want: &doc.Envelope{
				XMLNs: doc.EnvNamespace,
				SUM:   doc.SUM,
				SUM1:  doc.SUM1,
				Body: &doc.Body{
					VeriFactu: &doc.RegFactuSistemaFacturacion{
						Cabecera: &doc.Cabecera{
							Obligado: doc.Obligado{
								NombreRazon: "Test Company",
								NIF:         "B12345678",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid XML",
			data:    []byte(`invalid xml`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty document",
			data:    []byte{},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vf.ParseDocument(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDocument() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("should preserve parse whole doc", func(t *testing.T) {
		inv := test.LoadInvoice("inv-base.json")
		sw := &doc.Software{
			NombreRazon:              "My Software",
			NIF:                      "12345678A",
			NombreSistemaInformatico: "My Software",
			IdSistemaInformatico:     "A1",
			Version:                  "1.0",
			NumeroInstalacion:        "12345678A",
		}
		want, err := doc.NewVerifactu(inv, time.Now(), doc.IssuerRoleSupplier, sw, false)
		require.NoError(t, err)

		// Get the XML bytes from the reference document
		xmlData, err := want.Bytes()
		require.NoError(t, err)

		// Parse the XML back into a document
		got, err := vf.ParseDocument(xmlData)
		require.NoError(t, err)

		// Check that RegistroAlta is present and correctly structured
		require.NotNil(t, got.Body.VeriFactu.RegistroFactura)
		require.NotNil(t, got.Body.VeriFactu.RegistroFactura.RegistroAlta)
		assert.Equal(t, want.Body.VeriFactu.RegistroFactura.RegistroAlta.IDVersion,
			got.Body.VeriFactu.RegistroFactura.RegistroAlta.IDVersion)

	})
}

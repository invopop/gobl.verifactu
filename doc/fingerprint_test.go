package doc_test

import (
	"testing"

	"github.com/invopop/gobl.verifactu/doc"
)

var s = "S"

func TestFingerprintAlta(t *testing.T) {
	t.Run("Alta", func(t *testing.T) {
		tests := []struct {
			name     string
			alta     *doc.RegistroAlta
			prev     *doc.ChainData
			expected string
		}{
			{
				name: "Basic 1",
				alta: &doc.RegistroAlta{
					IDFactura: &doc.IDFactura{
						IDEmisorFactura:        "A28083806",
						NumSerieFactura:        "SAMPLE-001",
						FechaExpedicionFactura: "11-11-2024",
					},
					TipoFactura:              "F1",
					CuotaTotal:               378.0,
					ImporteTotal:             2178.0,
					FechaHoraHusoGenRegistro: "2024-11-20T19:00:55+01:00",
				},
				prev: &doc.ChainData{
					IDEmisorFactura:        "foo",
					NumSerieFactura:        "bar",
					FechaExpedicionFactura: "baz",
					Huella:                 "4B0A5C1D3F28E6A79B8C2D1E0F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8D9E0F1",
				},
				expected: "9F848AF7AECAA4C841654B37FD7119F4530B19141A2C3FF9968B5A229DEE21C2",
			},
			{
				name: "Basic 2",
				alta: &doc.RegistroAlta{
					IDFactura: &doc.IDFactura{
						IDEmisorFactura:        "A28083806",
						NumSerieFactura:        "SAMPLE-002",
						FechaExpedicionFactura: "12-11-2024",
					},
					TipoFactura:              "R3",
					CuotaTotal:               500.50,
					ImporteTotal:             2502.55,
					FechaHoraHusoGenRegistro: "2024-11-20T20:00:55+01:00",
				},
				prev: &doc.ChainData{
					IDEmisorFactura:        "foo",
					NumSerieFactura:        "bar",
					FechaExpedicionFactura: "baz",
					Huella:                 "4B0A5C1D3F28E6A79B8C2D1E0F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8D9E0F1",
				},
				expected: "14543C022CBD197F247F77A88F41E636A3B2569CE5787A8D6C8A781BF1B9D25E",
			},
			{
				name: "No Previous",
				alta: &doc.RegistroAlta{
					IDFactura: &doc.IDFactura{
						IDEmisorFactura:        "B08194359",
						NumSerieFactura:        "SAMPLE-003",
						FechaExpedicionFactura: "12-11-2024",
					},
					TipoFactura:              "F1",
					CuotaTotal:               500.0,
					ImporteTotal:             2500.0,
					Encadenamiento:           &doc.Encadenamiento{PrimerRegistro: &s},
					FechaHoraHusoGenRegistro: "2024-11-20T20:00:55+01:00",
				},
				prev:     nil,
				expected: "95619096010E699BB4B88AD2B42DC30BBD809A4B1ED2AE2904DFF86D064FCF29",
			},
			{
				name: "No Taxes",
				alta: &doc.RegistroAlta{
					IDFactura: &doc.IDFactura{
						IDEmisorFactura:        "B85905495",
						NumSerieFactura:        "SAMPLE-003",
						FechaExpedicionFactura: "15-11-2024",
					},
					TipoFactura:              "F1",
					CuotaTotal:               0.0,
					ImporteTotal:             1800.0,
					Encadenamiento:           &doc.Encadenamiento{RegistroAnterior: &doc.RegistroAnterior{Huella: "13EC0696104D1E529667184C6CDFC67D08036BCA4CD1B7887DE9C6F8F7EEC69C"}},
					FechaHoraHusoGenRegistro: "2024-11-21T17:59:41+01:00",
				},
				prev: &doc.ChainData{
					IDEmisorFactura:        "foo",
					NumSerieFactura:        "bar",
					FechaExpedicionFactura: "baz",
					Huella:                 "13EC0696104D1E529667184C6CDFC67D08036BCA4CD1B7887DE9C6F8F7EEC69C",
				},
				expected: "9F44F498EA51C0C50FEB026CCE86BDCCF852C898EE33336EFFE1BD6F132B506E",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := &doc.VeriFactu{
					RegistroFactura: &doc.RegistroFactura{
						RegistroAlta: tt.alta,
					},
				}

				err := d.Fingerprint(tt.prev)
				if err != nil {
					t.Errorf("fingerprintAlta() error = %v", err)
					return
				}

				if got := tt.alta.Huella; got != tt.expected {
					t.Errorf("fingerprint = %v, want %v", got, tt.expected)
				}
			})
		}
	})
}

func TestFingerprintAnulacion(t *testing.T) {
	t.Run("Anulacion", func(t *testing.T) {
		tests := []struct {
			name      string
			anulacion *doc.RegistroAnulacion
			prev      *doc.ChainData
			expected  string
		}{
			{
				name: "Basic 1",
				anulacion: &doc.RegistroAnulacion{
					IDFactura: &doc.IDFacturaAnulada{
						IDEmisorFactura:        "A28083806",
						NumSerieFactura:        "SAMPLE-001",
						FechaExpedicionFactura: "11-11-2024",
					},
					FechaHoraHusoGenRegistro: "2024-11-21T10:00:55+01:00",
					Encadenamiento:           &doc.Encadenamiento{RegistroAnterior: &doc.RegistroAnterior{Huella: "4B0A5C1D3F28E6A79B8C2D1E0F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8D9E0F1"}},
				},
				expected: "BAB9B4AE157321642F6AFD8030288B7E595129B29A00A69CEB308CEAA53BFBD7",
			},
			{
				name: "Basic 2",
				anulacion: &doc.RegistroAnulacion{
					IDFactura: &doc.IDFacturaAnulada{
						IDEmisorFactura:        "B08194359",
						NumSerieFactura:        "SAMPLE-002",
						FechaExpedicionFactura: "12-11-2024",
					},
					FechaHoraHusoGenRegistro: "2024-11-21T12:00:55+01:00",
					Encadenamiento:           &doc.Encadenamiento{RegistroAnterior: &doc.RegistroAnterior{Huella: "CBA051CBF59488B6978FA66E95ED4D0A84A97F5C0700EA952B923BD6E7C3FD7A"}},
				},
				expected: "548707E0984AA867CC173B24389E648DECDEE48A2674DA8CE8A3682EF8F119DD",
			},
			{
				name: "No Previous",
				anulacion: &doc.RegistroAnulacion{
					IDFactura: &doc.IDFacturaAnulada{
						IDEmisorFactura:        "A28083806",
						NumSerieFactura:        "SAMPLE-001",
						FechaExpedicionFactura: "11-11-2024",
					},
					FechaHoraHusoGenRegistro: "2024-11-21T10:00:55+01:00",
					Encadenamiento:           &doc.Encadenamiento{PrimerRegistro: &s},
				},
				expected: "CBA051CBF59488B6978FA66E95ED4D0A84A97F5C0700EA952B923BD6E7C3FD7A",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := &doc.VeriFactu{
					RegistroFactura: &doc.RegistroFactura{
						RegistroAnulacion: tt.anulacion,
					},
				}

				err := d.FingerprintCancel(tt.prev)
				if err != nil {
					t.Errorf("fingerprintAnulacion() error = %v", err)
					return
				}

				if got := tt.anulacion.Huella; got != tt.expected {
					t.Errorf("fingerprint = %v, want %v", got, tt.expected)
				}
			})
		}
	})
}

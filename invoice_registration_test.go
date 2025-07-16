package verifactu

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceRegistrationFingerprint(t *testing.T) {
	t.Run("Alta", func(t *testing.T) {
		tests := []struct {
			name     string
			alta     *InvoiceRegistration
			prev     *ChainData
			expected string
		}{
			{
				name: "Basic 1",
				alta: &InvoiceRegistration{
					IDFactura: &IDFactura{
						IDEmisorFactura:        "A28083806",
						NumSerieFactura:        "SAMPLE-001",
						FechaExpedicionFactura: "11-11-2024",
					},
					TipoFactura:              "F1",
					CuotaTotal:               num.MakeAmount(37800, 2),
					ImporteTotal:             num.MakeAmount(217800, 2),
					FechaHoraHusoGenRegistro: "2024-11-20T19:00:55+01:00",
				},
				prev: &ChainData{
					IDIssuer:    "A28083806",
					NumSeries:   "SAMPLE-000",
					IssueDate:   "10-11-2024",
					Fingerprint: "4B0A5C1D3F28E6A79B8C2D1E0F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8D9E0F1",
				},
				expected: "CDF61559670788B743FE269DF39622F01253932740BF6A254ABDABD622150E98",
			},
			{
				name: "Basic 2",
				alta: &InvoiceRegistration{
					IDFactura: &IDFactura{
						IDEmisorFactura:        "A28083806",
						NumSerieFactura:        "SAMPLE-002",
						FechaExpedicionFactura: "12-11-2024",
					},
					TipoFactura:              "R3",
					CuotaTotal:               num.MakeAmount(50050, 2),
					ImporteTotal:             num.MakeAmount(250255, 2),
					FechaHoraHusoGenRegistro: "2024-11-20T20:00:55+01:00",
				},
				prev: &ChainData{
					IDIssuer:    "A28083806",
					NumSeries:   "SAMPLE-001",
					IssueDate:   "11-11-2024",
					Fingerprint: "4B0A5C1D3F28E6A79B8C2D1E0F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8D9E0F1",
				},
				expected: "7482C035254A112DE6359FF01A43ABBA4E7318D434A243D79F2915768AB07E06",
			},
			{
				name: "No Previous",
				alta: &InvoiceRegistration{
					IDFactura: &IDFactura{
						IDEmisorFactura:        "B08194359",
						NumSerieFactura:        "SAMPLE-003",
						FechaExpedicionFactura: "12-11-2024",
					},
					TipoFactura:              "F1",
					CuotaTotal:               num.MakeAmount(50000, 2),
					ImporteTotal:             num.MakeAmount(250000, 2),
					Encadenamiento:           &Encadenamiento{PrimerRegistro: "S"},
					FechaHoraHusoGenRegistro: "2024-11-20T20:00:55+01:00",
				},
				prev:     nil,
				expected: "8077BD9EE0A9EA26AB0FB86762B20379264B7DC8F8DBE2972D7FEBC5D6B0B915",
			},
			{
				name: "No Taxes",
				alta: &InvoiceRegistration{
					IDFactura: &IDFactura{
						IDEmisorFactura:        "B85905495",
						NumSerieFactura:        "SAMPLE-003",
						FechaExpedicionFactura: "15-11-2024",
					},
					TipoFactura:              "F1",
					CuotaTotal:               num.MakeAmount(0, 2),
					ImporteTotal:             num.MakeAmount(180000, 2),
					Encadenamiento:           &Encadenamiento{RegistroAnterior: &RegistroAnterior{Huella: "13EC0696104D1E529667184C6CDFC67D08036BCA4CD1B7887DE9C6F8F7EEC69C"}},
					FechaHoraHusoGenRegistro: "2024-11-21T17:59:41+01:00",
				},
				prev: &ChainData{
					IDIssuer:    "A28083806",
					NumSeries:   "SAMPLE-002",
					IssueDate:   "11-11-2024",
					Fingerprint: "13EC0696104D1E529667184C6CDFC67D08036BCA4CD1B7887DE9C6F8F7EEC69C",
				},
				expected: "C95799034D369BB42BABCB996AA962EB0FC9B66D58F0E35F610113BF4B15A3FD",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.alta.fingerprint(tt.prev)
				assert.Equal(t, tt.expected, tt.alta.Huella)
			})
		}
	})
}

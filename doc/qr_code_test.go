package doc

import (
	"testing"
)

func TestGenerateCodes(t *testing.T) {
	tests := []struct {
		name     string
		doc      *VeriFactu
		expected string
	}{
		{
			name: "valid codes generation",
			doc: &VeriFactu{
				RegistroFactura: &RegistroFactura{
					RegistroAlta: &RegistroAlta{
						IDFactura: &IDFactura{
							IDEmisorFactura:        "89890001K",
							NumSerieFactura:        "12345678-G33",
							FechaExpedicionFactura: "01-09-2024",
						},
						ImporteTotal: 241.4,
					},
				},
			},
			expected: "https://prewww2.aeat.es/wlpl/TIKE-CONT/ValidarQR?nif=89890001K&numserie=12345678-G33&fecha=01-09-2024&importe=241.4",
		},
		{
			name: "empty fields",
			doc: &VeriFactu{
				RegistroFactura: &RegistroFactura{
					RegistroAlta: &RegistroAlta{
						IDFactura: &IDFactura{
							IDEmisorFactura:        "",
							NumSerieFactura:        "",
							FechaExpedicionFactura: "",
						},
						ImporteTotal: 0,
					},
				},
			},
			expected: "https://prewww2.aeat.es/wlpl/TIKE-CONT/ValidarQR?nif=&numserie=&fecha=&importe=0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.doc.generateURL()
			if got != tt.expected {
				t.Errorf("generateURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGenerateURLCodeAlta(t *testing.T) {
	tests := []struct {
		name     string
		doc      *VeriFactu
		expected string
	}{
		{
			name: "valid URL generation",
			doc: &VeriFactu{
				RegistroFactura: &RegistroFactura{
					RegistroAlta: &RegistroAlta{
						IDFactura: &IDFactura{
							IDEmisorFactura:        "89890001K",
							NumSerieFactura:        "12345678-G33",
							FechaExpedicionFactura: "01-09-2024",
						},
						ImporteTotal: 241.4,
					},
				},
			},
			expected: "https://prewww2.aeat.es/wlpl/TIKE-CONT/ValidarQR?nif=89890001K&numserie=12345678-G33&fecha=01-09-2024&importe=241.4",
		},
		{
			name: "URL with special characters",
			doc: &VeriFactu{
				RegistroFactura: &RegistroFactura{
					RegistroAlta: &RegistroAlta{
						IDFactura: &IDFactura{
							IDEmisorFactura:        "A12 345&67",
							NumSerieFactura:        "SERIE/2023",
							FechaExpedicionFactura: "01-09-2024",
						},
						ImporteTotal: 1234.56,
					},
				},
			},
			expected: "https://prewww2.aeat.es/wlpl/TIKE-CONT/ValidarQR?nif=A12+345%2667&numserie=SERIE%2F2023&fecha=01-09-2024&importe=1234.56",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.doc.generateURL()
			if got != tt.expected {
				t.Errorf("generateURLCodeAlta() = %v, want %v", got, tt.expected)
			}
		})
	}
}

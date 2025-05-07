package verifactu

import (
	"testing"
)

func TestGenerateCodes(t *testing.T) {
	tests := []struct {
		name     string
		doc      *InvoiceRegistration
		expected string
	}{
		{
			name: "valid codes generation",
			doc: &InvoiceRegistration{
				IDFactura: &IDFactura{
					IDEmisorFactura:        "89890001K",
					NumSerieFactura:        "12345678-G33",
					FechaExpedicionFactura: "01-09-2024",
				},
				ImporteTotal: "241.40",
			},
			expected: "https://prewww2.aeat.es/wlpl/TIKE-CONT/ValidarQR?nif=89890001K&numserie=12345678-G33&fecha=01-09-2024&importe=241.40",
		},
		{
			name: "empty fields",
			doc: &InvoiceRegistration{
				IDFactura: &IDFactura{
					IDEmisorFactura:        "",
					NumSerieFactura:        "",
					FechaExpedicionFactura: "",
				},
				ImporteTotal: "0.00",
			},
			expected: "https://prewww2.aeat.es/wlpl/TIKE-CONT/ValidarQR?nif=&numserie=&fecha=&importe=0.00",
		},
		{
			name: "URL with special characters",
			doc: &InvoiceRegistration{
				IDFactura: &IDFactura{
					IDEmisorFactura:        "A12 345&67",
					NumSerieFactura:        "SERIE/2023",
					FechaExpedicionFactura: "01-09-2024",
				},
				ImporteTotal: "1234.56",
			},
			expected: "https://prewww2.aeat.es/wlpl/TIKE-CONT/ValidarQR?nif=A12+345%2667&numserie=SERIE%2F2023&fecha=01-09-2024&importe=1234.56",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.doc.generateURL(false)
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

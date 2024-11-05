package doc

import "github.com/invopop/gobl/bill"

type Software struct {
	NombreRazon                 string
	NIF                         string
	IdSistemaInformatico        string
	NombreSistemaInformatico    string
	NumeroInstalacion           string
	TipoUsoPosibleSoloVerifactu string
	TipoUsoPosibleMultiOT       string
	IndicadorMultiplesOT        string
	Version                     string
}

func newSoftware(inv *bill.Invoice) *Software {
	software := &Software{
		NombreRazon:          "Invopop SL",
		NIF:                  inv.Supplier.TaxID.Code.String(),
		IDOtro:               "04",
		IDSistemaInformatico: "F1",
		Version:              "1.0",
		NumeroInstalacion:    "1",
	}
	return software
}

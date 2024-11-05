package doc

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

func newSoftware() *Software {
	software := &Software{
		NombreRazon: "xxxxxxxx",
		NIF:         "0123456789",
		// IDOtro:               "04",
		// IDSistemaInformatico: "F1",
		Version:           "1.0",
		NumeroInstalacion: "1",
	}
	return software
}

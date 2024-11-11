package doc

type Software struct {
	NombreRazon string
	NIF         string
	// IDOtro                      string
	// NombreSistemaInformatico    string
	IdSistemaInformatico string
	Version              string
	NumeroInstalacion    string
	// TipoUsoPosibleSoloVerifactu string
	// TipoUsoPosibleMultiOT       string
	// IndicadorMultiplesOT        string
}

func newSoftware() *Software {
	software := &Software{
		NombreRazon: "xxxxxxxx",
		NIF:         "0123456789",
		// IDOtro:               "04",
		// IDSistemaInformatico: "F1",
		Version:           CurrentVersion,
		NumeroInstalacion: "1",
	}
	return software
}

package doc

type VeriFactu struct {
	Cabecera        *Cabecera
	RegistroFactura []*RegistroFactura
}

type RegistroFactura struct {
	RegistroAlta      *RegistroAlta
	RegistroAnulacion *RegistroAnulacion
}

type RegistroAnulacion struct {
}

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

package doc

type RegistroAlta struct {
	IDVersion                           string                `xml:"IDVersion"`
	IDFactura                           IDFactura             `xml:"IDFactura"`
	RefExterna                          string                `xml:"RefExterna,omitempty"`
	NombreRazonEmisor                   string                `xml:"NombreRazonEmisor"`
	Subsanacion                         string                `xml:"Subsanacion,omitempty"`
	RechazoPrevio                       string                `xml:"RechazoPrevio,omitempty"`
	TipoFactura                         string                `xml:"TipoFactura"`
	TipoRectificativa                   string                `xml:"TipoRectificativa,omitempty"`
	FacturasRectificadas                []*FacturaRectificada `xml:"FacturasRectificadas>FacturaRectificada,omitempty"`
	FacturasSustituidas                 []*FacturaSustituida  `xml:"FacturasSustituidas>FacturaSustituida,omitempty"`
	ImporteRectificacion                ImporteRectificacion  `xml:"ImporteRectificacion,omitempty"`
	FechaOperacion                      string                `xml:"FechaOperacion"`
	DescripcionOperacion                string                `xml:"DescripcionOperacion"`
	FacturaSimplificadaArt7273          string                `xml:"FacturaSimplificadaArt7273,omitempty"`
	FacturaSinIdentifDestinatarioArt61d string                `xml:"FacturaSinIdentifDestinatarioArt61d,omitempty"`
	Macrodato                           string                `xml:"Macrodato,omitempty"`
	EmitidaPorTerceroODestinatario      string                `xml:"EmitidaPorTerceroODestinatario,omitempty"`
	Tercero                             Tercero               `xml:"Tercero,omitempty"`
	Destinatarios                       []*Destinatario       `xml:"Destinatarios>Destinatario,omitempty"`
	Cupon                               string                `xml:"Cupon,omitempty"`
	Desglose                            Desglose              `xml:"Desglose"`
	CuotaTotal                          float64               `xml:"CuotaTotal"`
	ImporteTotal                        float64               `xml:"ImporteTotal"`
	Encadenamiento                      Encadenamiento        `xml:"Encadenamiento"`
	SistemaInformatico                  Software              `xml:"SistemaInformatico"`
	FechaHoraHusoGenRegistro            string                `xml:"FechaHoraHusoGenRegistro"`
	NumRegistroAcuerdoFacturacion       string                `xml:"NumRegistroAcuerdoFacturacion,omitempty"`
	IdAcuerdoSistemaInformatico         string                `xml:"IdAcuerdoSistemaInformatico,omitempty"`
	TipoHuella                          string                `xml:"TipoHuella"`
	Huella                              string                `xml:"Huella"`
	Signature                           string                `xml:"Signature"`
}

type IDFactura struct {
	IDEmisorFactura        string `xml:"IDEmisorFactura"`
	NumSerieFactura        string `xml:"NumSerieFactura"`
	FechaExpedicionFactura string `xml:"FechaExpedicionFactura"`
}

type FacturaRectificada struct {
	IDFactura IDFactura `xml:"IDFactura"`
}

type FacturaSustituida struct {
	IDFactura IDFactura `xml:"IDFactura"`
}

type ImporteRectificacion struct {
	BaseRectificada         float64 `xml:"BaseRectificada"`
	CuotaRectificada        float64 `xml:"CuotaRectificada"`
	CuotaRecargoRectificado float64 `xml:"CuotaRecargoRectificado"`
}

type Tercero struct {
	Nif         string `xml:"Nif,omitempty"`
	NombreRazon string `xml:"NombreRazon"`
	IDOtro      string `xml:"IDOtro,omitempty"`
}

type Destinatario struct {
	IDDestinatario IDDestinatario `xml:"IDDestinatario"`
}

type IDDestinatario struct {
	NIF         string `xml:"NIF,omitempty"`
	NombreRazon string `xml:"NombreRazon"`
	IDOtro      IDOtro `xml:"IDOtro,omitempty"`
}

type IDOtro struct {
	CodigoPais string `xml:"CodigoPais"`
	IDType     string `xml:"IDType"`
	ID         string `xml:"ID"`
}

type Desglose struct {
	DetalleDesglose []*DetalleDesglose `xml:"DetalleDesglose"`
}

type DetalleDesglose struct {
	Impuesto                      string `xml:"Impuesto"`
	ClaveRegimen                  string `xml:"ClaveRegimen"`
	CalificacionOperacion         string `xml:"CalificacionOperacion,omitempty"`
	OperacionExenta               string `xml:"OperacionExenta,omitempty"`
	TipoImpositivo                string `xml:"TipoImpositivo,omitempty"`
	BaseImponibleOImporteNoSujeto string `xml:"BaseImponibleOImporteNoSujeto"`
	BaseImponibleACoste           string `xml:"BaseImponibleACoste,omitempty"`
	CuotaRepercutida              string `xml:"CuotaRepercutida,omitempty"`
	TipoRecargoEquivalencia       string `xml:"TipoRecargoEquivalencia,omitempty"`
	CuotaRecargoEquivalencia      string `xml:"CuotaRecargoEquivalencia,omitempty"`
}

type Encadenamiento struct {
	PrimerRegistro   string           `xml:"PrimerRegistro"`
	RegistroAnterior RegistroAnterior `xml:"RegistroAnterior,omitempty"`
}

type RegistroAnterior struct {
	IDEmisorFactura        string `xml:"IDEmisorFactura"`
	NumSerieFactura        string `xml:"NumSerieFactura"`
	FechaExpedicionFactura string `xml:"FechaExpedicionFactura"`
	Huella                 string `xml:"Huella"`
}

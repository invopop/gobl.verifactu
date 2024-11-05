package doc

// "github.com/invopop/gobl/pkg/xmldsig"

const (
	SUM  = "https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroLR.xsd"
	SUM1 = "https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroInformacion.xsd"
)

type Cabecera struct {
	Obligado              Obligado               `xml:"sum1:Obligado"`
	Representante         *Obligado              `xml:"sum1:Representante,omitempty"`
	RemisionVoluntaria    *RemisionVoluntaria    `xml:"sum1:RemisionVoluntaria,omitempty"`
	RemisionRequerimiento *RemisionRequerimiento `xml:"sum1:RemisionRequerimiento,omitempty"`
}

type Obligado struct {
	NombreRazon string `xml:"sum1:NombreRazon"`
	NIF         string `xml:"sum1:NIF"`
}

type RemisionVoluntaria struct {
	FechaFinVerifactu string `xml:"sum1:FechaFinVerifactu"`
	Incidencia        string `xml:"sum1:Incidencia"`
}

type RemisionRequerimiento struct {
	RefRequerimiento string `xml:"sum1:RefRequerimiento"`
	FinRequerimiento string `xml:"sum1:FinRequerimiento"`
}

type RegistroAlta struct {
	IDVersion                           string                `xml:"sum1:IDVersion"`
	IDFactura                           *IDFactura            `xml:"sum1:IDFactura"`
	RefExterna                          string                `xml:"sum1:RefExterna,omitempty"`
	NombreRazonEmisor                   string                `xml:"sum1:NombreRazonEmisor"`
	Subsanacion                         string                `xml:"sum1:Subsanacion,omitempty"`
	RechazoPrevio                       string                `xml:"sum1:RechazoPrevio,omitempty"`
	TipoFactura                         string                `xml:"sum1:TipoFactura"`
	TipoRectificativa                   string                `xml:"sum1:TipoRectificativa,omitempty"`
	FacturasRectificadas                []*FacturaRectificada `xml:"sum1:FacturasRectificadas>sum1:FacturaRectificada,omitempty"`
	FacturasSustituidas                 []*FacturaSustituida  `xml:"sum1:FacturasSustituidas>sum1:FacturaSustituida,omitempty"`
	ImporteRectificacion                *ImporteRectificacion `xml:"sum1:ImporteRectificacion,omitempty"`
	FechaOperacion                      string                `xml:"sum1:FechaOperacion"`
	DescripcionOperacion                string                `xml:"sum1:DescripcionOperacion"`
	FacturaSimplificadaArt7273          string                `xml:"sum1:FacturaSimplificadaArt7273,omitempty"`
	FacturaSinIdentifDestinatarioArt61d string                `xml:"sum1:FacturaSinIdentifDestinatarioArt61d,omitempty"`
	Macrodato                           string                `xml:"sum1:Macrodato,omitempty"`
	EmitidaPorTerceroODestinatario      string                `xml:"sum1:EmitidaPorTerceroODestinatario,omitempty"`
	Tercero                             *Tercero              `xml:"sum1:Tercero,omitempty"`
	Destinatarios                       []*Destinatario       `xml:"sum1:Destinatarios>sum1:Destinatario,omitempty"`
	Cupon                               string                `xml:"sum1:Cupon,omitempty"`
	Desglose                            *Desglose             `xml:"sum1:Desglose"`
	CuotaTotal                          string                `xml:"sum1:CuotaTotal"`
	ImporteTotal                        string                `xml:"sum1:ImporteTotal"`
	Encadenamiento                      *Encadenamiento       `xml:"sum1:Encadenamiento"`
	SistemaInformatico                  *Software             `xml:"sum1:SistemaInformatico"`
	FechaHoraHusoGenRegistro            string                `xml:"sum1:FechaHoraHusoGenRegistro"`
	NumRegistroAcuerdoFacturacion       string                `xml:"sum1:NumRegistroAcuerdoFacturacion,omitempty"`
	IdAcuerdoSistemaInformatico         string                `xml:"sum1:IdAcuerdoSistemaInformatico,omitempty"`
	TipoHuella                          string                `xml:"sum1:TipoHuella"`
	Huella                              string                `xml:"sum1:Huella"`
	// Signature                           *xmldsig.Signature   `xml:"sum1:Signature,omitempty"`
}

type IDFactura struct {
	IDEmisorFactura        string `xml:"sum1:IDEmisorFactura"`
	NumSerieFactura        string `xml:"sum1:NumSerieFactura"`
	FechaExpedicionFactura string `xml:"sum1:FechaExpedicionFactura"`
}

type FacturaRectificada struct {
	IDFactura IDFactura `xml:"sum1:IDFactura"`
}

type FacturaSustituida struct {
	IDFactura IDFactura `xml:"sum1:IDFactura"`
}

type ImporteRectificacion struct {
	BaseRectificada         string `xml:"sum1:BaseRectificada"`
	CuotaRectificada        string `xml:"sum1:CuotaRectificada"`
	CuotaRecargoRectificado string `xml:"sum1:CuotaRecargoRectificado"`
}

type Tercero struct {
	Nif         string `xml:"sum1:Nif,omitempty"`
	NombreRazon string `xml:"sum1:NombreRazon"`
	IDOtro      string `xml:"sum1:IDOtro,omitempty"`
}

type Destinatario struct {
	IDDestinatario IDDestinatario `xml:"sum1:IDDestinatario"`
}

type IDDestinatario struct {
	NIF         string `xml:"sum1:NIF,omitempty"`
	NombreRazon string `xml:"sum1:NombreRazon"`
	IDOtro      IDOtro `xml:"sum1:IDOtro,omitempty"`
}

type IDOtro struct {
	CodigoPais string `xml:"sum1:CodigoPais"`
	IDType     string `xml:"sum1:IDType"`
	ID         string `xml:"sum1:ID"`
}

type Desglose struct {
	DetalleDesglose []*DetalleDesglose `xml:"sum1:DetalleDesglose"`
}

type DetalleDesglose struct {
	Impuesto                      string `xml:"sum1:Impuesto"`
	ClaveRegimen                  string `xml:"sum1:ClaveRegimen"`
	CalificacionOperacion         string `xml:"sum1:CalificacionOperacion,omitempty"`
	OperacionExenta               string `xml:"sum1:OperacionExenta,omitempty"`
	TipoImpositivo                string `xml:"sum1:TipoImpositivo,omitempty"`
	BaseImponibleOImporteNoSujeto string `xml:"sum1:BaseImponibleOImporteNoSujeto"`
	BaseImponibleACoste           string `xml:"sum1:BaseImponibleACoste,omitempty"`
	CuotaRepercutida              string `xml:"sum1:CuotaRepercutida,omitempty"`
	TipoRecargoEquivalencia       string `xml:"sum1:TipoRecargoEquivalencia,omitempty"`
	CuotaRecargoEquivalencia      string `xml:"sum1:CuotaRecargoEquivalencia,omitempty"`
}

type Encadenamiento struct {
	PrimerRegistro   string           `xml:"sum1:PrimerRegistro"`
	RegistroAnterior RegistroAnterior `xml:"sum1:RegistroAnterior,omitempty"`
}

type RegistroAnterior struct {
	IDEmisorFactura        string `xml:"sum1:IDEmisorFactura"`
	NumSerieFactura        string `xml:"sum1:NumSerieFactura"`
	FechaExpedicionFactura string `xml:"sum1:FechaExpedicionFactura"`
	Huella                 string `xml:"sum1:Huella"`
}

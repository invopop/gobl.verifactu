package doc

import "encoding/xml"

// SUM is the namespace for the main VeriFactu schema
const (
	SUM          = "https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroLR.xsd"
	SUM1         = "https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroInformacion.xsd"
	EnvNamespace = "http://schemas.xmlsoap.org/soap/envelope/"
)

//	xmlns:sf="https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroInformacion.xsd"
//	xmlns:sfLR="https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroLR.xsd"
//	xmlns:sfR="https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/RespuestaSuministro.xsd"
//
// Envelope is the SOAP envelope wrapper
type Envelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	XMLNs   string   `xml:"xmlns:soapenv,attr"`
	SUM     string   `xml:"xmlns:sum,attr"`
	SUM1    string   `xml:"xmlns:sum1,attr"`
	Body    struct {
		XMLName   xml.Name   `xml:"soapenv:Body"`
		VeriFactu *VeriFactu `xml:"sum:RegFactuSistemaFacturacion"`
	}
}

// VeriFactu represents the root element of a VeriFactu document
type VeriFactu struct {
	XMLName         xml.Name         `xml:"sum:RegFactuSistemaFacturacion"`
	Cabecera        *Cabecera        `xml:"sum:Cabecera"`
	RegistroFactura *RegistroFactura `xml:"sum:RegistroFactura"`
}

// RegistroFactura contains either an invoice registration or cancellation
type RegistroFactura struct {
	RegistroAlta      *RegistroAlta      `xml:"sum1:RegistroAlta,omitempty"`
	RegistroAnulacion *RegistroAnulacion `xml:"sum1:RegistroAnulacion,omitempty"`
}

// Cabecera contains the header information for a VeriFactu document
type Cabecera struct {
	Obligado              Obligado               `xml:"sum1:ObligadoEmision"`
	Representante         *Obligado              `xml:"sum1:Representante,omitempty"`
	RemisionVoluntaria    *RemisionVoluntaria    `xml:"sum1:RemisionVoluntaria,omitempty"`
	RemisionRequerimiento *RemisionRequerimiento `xml:"sum1:RemisionRequerimiento,omitempty"`
}

// Obligado represents an obligated party in the document
type Obligado struct {
	NombreRazon string `xml:"sum1:NombreRazon"`
	NIF         string `xml:"sum1:NIF"`
}

// RemisionVoluntaria contains voluntary submission details
type RemisionVoluntaria struct {
	FechaFinVerifactu string `xml:"sum1:FechaFinVerifactu"`
	Incidencia        string `xml:"sum1:Incidencia"`
}

// RemisionRequerimiento contains requirement submission details
type RemisionRequerimiento struct {
	RefRequerimiento string `xml:"sum1:RefRequerimiento"`
	FinRequerimiento string `xml:"sum1:FinRequerimiento"`
}

// RegistroAlta contains the details of an invoice registration
type RegistroAlta struct {
	IDVersion                           string                `xml:"sum1:IDVersion"`
	IDFactura                           *IDFactura            `xml:"sum1:IDFactura"`
	RefExterna                          string                `xml:"sum1:RefExterna,omitempty"`
	NombreRazonEmisor                   string                `xml:"sum1:NombreRazonEmisor"`
	Subsanacion                         string                `xml:"sum1:Subsanacion,omitempty"`
	RechazoPrevio                       string                `xml:"sum1:RechazoPrevio,omitempty"`
	TipoFactura                         string                `xml:"sum1:TipoFactura"`
	TipoRectificativa                   string                `xml:"sum1:TipoRectificativa,omitempty"`
	FacturasRectificadas                []*FacturaRectificada `xml:"sum1:FacturasRectificadas,omitempty"`
	FacturasSustituidas                 []*FacturaSustituida  `xml:"sum1:FacturasSustituidas,omitempty"`
	ImporteRectificacion                *ImporteRectificacion `xml:"sum1:ImporteRectificacion,omitempty"`
	FechaOperacion                      string                `xml:"sum1:FechaOperacion"`
	DescripcionOperacion                string                `xml:"sum1:DescripcionOperacion"`
	FacturaSimplificadaArt7273          string                `xml:"sum1:FacturaSimplificadaArt7273,omitempty"`
	FacturaSinIdentifDestinatarioArt61d string                `xml:"sum1:FacturaSinIdentifDestinatarioArt61d,omitempty"`
	Macrodato                           string                `xml:"sum1:Macrodato,omitempty"`
	EmitidaPorTerceroODestinatario      string                `xml:"sum1:EmitidaPorTerceroODestinatario,omitempty"`
	Tercero                             *Party                `xml:"sum1:Tercero,omitempty"`
	Destinatarios                       []*Destinatario       `xml:"sum1:Destinatarios,omitempty"`
	Cupon                               string                `xml:"sum1:Cupon,omitempty"`
	Desglose                            *Desglose             `xml:"sum1:Desglose"`
	CuotaTotal                          float64               `xml:"sum1:CuotaTotal"`
	ImporteTotal                        float64               `xml:"sum1:ImporteTotal"`
	Encadenamiento                      *Encadenamiento       `xml:"sum1:Encadenamiento"`
	SistemaInformatico                  *Software             `xml:"sum1:SistemaInformatico"`
	FechaHoraHusoGenRegistro            string                `xml:"sum1:FechaHoraHusoGenRegistro"`
	NumRegistroAcuerdoFacturacion       string                `xml:"sum1:NumRegistroAcuerdoFacturacion,omitempty"`
	IdAcuerdoSistemaInformatico         string                `xml:"sum1:IdAcuerdoSistemaInformatico,omitempty"` //nolint:revive
	TipoHuella                          string                `xml:"sum1:TipoHuella"`
	Huella                              string                `xml:"sum1:Huella"`
	// Signature                           *xmldsig.Signature   `xml:"sum1:Signature,omitempty"`
}

// RegistroAnulacion contains the details of an invoice cancellation
type RegistroAnulacion struct {
	IDVersion                string          `xml:"IDVersion"`
	IDFactura                *IDFactura      `xml:"IDFactura"`
	RefExterna               string          `xml:"RefExterna,omitempty"`
	SinRegistroPrevio        string          `xml:"SinRegistroPrevio"`
	RechazoPrevio            string          `xml:"RechazoPrevio,omitempty"`
	GeneradoPor              string          `xml:"GeneradoPor"`
	Generador                *Party          `xml:"Generador"`
	Encadenamiento           *Encadenamiento `xml:"Encadenamiento"`
	SistemaInformatico       *Software       `xml:"SistemaInformatico"`
	FechaHoraHusoGenRegistro string          `xml:"FechaHoraHusoGenRegistro"`
	TipoHuella               string          `xml:"TipoHuella"`
	Huella                   string          `xml:"Huella"`
	Signature                string          `xml:"Signature"`
}

// IDFactura contains the identifying information for an invoice
type IDFactura struct {
	IDEmisorFactura        string `xml:"sum1:IDEmisorFactura"`
	NumSerieFactura        string `xml:"sum1:NumSerieFactura"`
	FechaExpedicionFactura string `xml:"sum1:FechaExpedicionFactura"`
}

// FacturaRectificada represents a rectified invoice
type FacturaRectificada struct {
	IDFactura IDFactura `xml:"sum1:IDFactura"`
}

// FacturaSustituida represents a substituted invoice
type FacturaSustituida struct {
	IDFactura IDFactura `xml:"sum1:IDFactura"`
}

// ImporteRectificacion contains rectification amounts
type ImporteRectificacion struct {
	BaseRectificada         string `xml:"sum1:BaseRectificada"`
	CuotaRectificada        string `xml:"sum1:CuotaRectificada"`
	CuotaRecargoRectificado string `xml:"sum1:CuotaRecargoRectificado"`
}

// Party represents a in the document, covering fields Generador, Tercero and IDDestinatario
type Party struct {
	NombreRazon string  `xml:"sum1:NombreRazon"`
	NIF         string  `xml:"sum1:NIF,omitempty"`
	IDOtro      *IDOtro `xml:"sum1:IDOtro,omitempty"`
}

// Destinatario represents a recipient in the document
type Destinatario struct {
	IDDestinatario *Party `xml:"sum1:IDDestinatario"`
}

// IDOtro contains alternative identifying information
type IDOtro struct {
	CodigoPais string `xml:"sum1:CodigoPais"`
	IDType     string `xml:"sum1:IDType"`
	ID         string `xml:"sum1:ID"`
}

// Desglose contains the breakdown details
type Desglose struct {
	DetalleDesglose []*DetalleDesglose `xml:"sum1:DetalleDesglose"`
}

// DetalleDesglose contains detailed breakdown information
type DetalleDesglose struct {
	Impuesto                      string  `xml:"sum1:Impuesto,omitempty"`
	ClaveRegimen                  string  `xml:"sum1:ClaveRegimen,omitempty"`
	CalificacionOperacion         string  `xml:"sum1:CalificacionOperacion,omitempty"`
	OperacionExenta               string  `xml:"sum1:OperacionExenta,omitempty"`
	TipoImpositivo                float64 `xml:"sum1:TipoImpositivo,omitempty"`
	BaseImponibleOImporteNoSujeto float64 `xml:"sum1:BaseImponibleOimporteNoSujeto"`
	BaseImponibleACoste           float64 `xml:"sum1:BaseImponibleACoste,omitempty"`
	CuotaRepercutida              float64 `xml:"sum1:CuotaRepercutida,omitempty"`
	TipoRecargoEquivalencia       float64 `xml:"sum1:TipoRecargoEquivalencia,omitempty"`
	CuotaRecargoEquivalencia      float64 `xml:"sum1:CuotaRecargoEquivalencia,omitempty"`
}

// Encadenamiento contains chaining information between documents
type Encadenamiento struct {
	PrimerRegistro   *string          `xml:"sum1:PrimerRegistro,omitempty"`
	RegistroAnterior RegistroAnterior `xml:"sum1:RegistroAnterior,omitempty"`
}

// RegistroAnterior contains information about the previous registration
type RegistroAnterior struct {
	IDEmisorFactura        string `xml:"sum1:IDEmisorFactura"`
	NumSerieFactura        string `xml:"sum1:NumSerieFactura"`
	FechaExpedicionFactura string `xml:"sum1:FechaExpedicionFactura"`
	Huella                 string `xml:"sum1:Huella"`
}

// Software contains the details about the software that is using this library to
// generate VeriFactu documents. These details are included in the final
// document.
type Software struct {
	NombreRazon                 string `xml:"sum1:NombreRazon"`
	NIF                         string `xml:"sum1:NIF"`
	NombreSistemaInformatico    string `xml:"sum1:NombreSistemaInformatico"`
	IdSistemaInformatico        string `xml:"sum1:IdSistemaInformatico"` //nolint:revive
	Version                     string `xml:"sum1:Version"`
	NumeroInstalacion           string `xml:"sum1:NumeroInstalacion"`
	TipoUsoPosibleSoloVerifactu string `xml:"sum1:TipoUsoPosibleSoloVerifactu,omitempty"`
	TipoUsoPosibleMultiOT       string `xml:"sum1:TipoUsoPosibleMultiOT,omitempty"`
	IndicadorMultiplesOT        string `xml:"sum1:IndicadorMultiplesOT,omitempty"`
}

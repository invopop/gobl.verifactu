package doc

import (
	"encoding/xml"

	"github.com/invopop/gobl/num"
)

// SUM is the namespace for the main VeriFactu schema
const (
	SUM          = "https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroLR.xsd"
	SUM1         = "https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroInformacion.xsd"
	EnvNamespace = "http://schemas.xmlsoap.org/soap/envelope/"
)

// Envelope is the SOAP envelope wrapper
type Envelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	XMLNs   string   `xml:"xmlns:soapenv,attr"`
	SUM     string   `xml:"xmlns:sum,attr"`
	SUM1    string   `xml:"xmlns:sum1,attr"`
	Body    *Body    `xml:"soapenv:Body"`
}

// Body is the body of the SOAP envelope
type Body struct {
	VeriFactu *RegFactuSistemaFacturacion `xml:"sum:RegFactuSistemaFacturacion"`
}

// RegFactuSistemaFacturacion represents the root element of a RegFactuSistemaFacturacion document
type RegFactuSistemaFacturacion struct {
	// XMLName         xml.Name         `xml:"sum:RegFactuSistemaFacturacion"`
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
	FechaFinVerifactu string `xml:"sum1:FechaFinVerifactu,omitempty"`
	Incidencia        string `xml:"sum1:Incidencia,omitempty"`
}

// RemisionRequerimiento contains requirement submission details
type RemisionRequerimiento struct {
	RefRequerimiento string `xml:"sum1:RefRequerimiento"`
	FinRequerimiento string `xml:"sum1:FinRequerimiento,omitempty"`
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
	FechaOperacion                      string                `xml:"sum1:FechaOperacion,omitempty"`
	DescripcionOperacion                string                `xml:"sum1:DescripcionOperacion"`
	FacturaSimplificadaArt7273          string                `xml:"sum1:FacturaSimplificadaArt7273,omitempty"`
	FacturaSinIdentifDestinatarioArt61d string                `xml:"sum1:FacturaSinIdentifDestinatarioArt61d,omitempty"`
	Macrodato                           string                `xml:"sum1:Macrodato,omitempty"`
	EmitidaPorTerceroODestinatario      string                `xml:"sum1:EmitidaPorTerceroODestinatario,omitempty"`
	Tercero                             *Party                `xml:"sum1:Tercero,omitempty"`
	Destinatarios                       []*Destinatario       `xml:"sum1:Destinatarios,omitempty"`
	Cupon                               string                `xml:"sum1:Cupon,omitempty"`
	Desglose                            *Desglose             `xml:"sum1:Desglose"`
	CuotaTotal                          string                `xml:"sum1:CuotaTotal"`
	ImporteTotal                        string                `xml:"sum1:ImporteTotal"`
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
	IDVersion                string            `xml:"sum1:IDVersion"`
	IDFactura                *IDFacturaAnulada `xml:"sum1:IDFactura"`
	RefExterna               string            `xml:"sum1:RefExterna,omitempty"`
	SinRegistroPrevio        string            `xml:"sum1:SinRegistroPrevio,omitempty"`
	RechazoPrevio            string            `xml:"sum1:RechazoPrevio,omitempty"`
	GeneradoPor              string            `xml:"sum1:GeneradoPor,omitempty"`
	Generador                *Party            `xml:"sum1:Generador,omitempty"`
	Encadenamiento           *Encadenamiento   `xml:"sum1:Encadenamiento"`
	SistemaInformatico       *Software         `xml:"sum1:SistemaInformatico"`
	FechaHoraHusoGenRegistro string            `xml:"sum1:FechaHoraHusoGenRegistro"`
	TipoHuella               string            `xml:"sum1:TipoHuella"`
	Huella                   string            `xml:"sum1:Huella"`
	// Signature               *xmldsig.Signature            `xml:"sum1:Signature"`
}

// IDFactura contains the identifying information for an invoice
type IDFactura struct {
	IDEmisorFactura        string `xml:"sum1:IDEmisorFactura"`
	NumSerieFactura        string `xml:"sum1:NumSerieFactura"`
	FechaExpedicionFactura string `xml:"sum1:FechaExpedicionFactura"`
}

// IDFacturaAnulada contains the identifying information for an invoice
type IDFacturaAnulada struct {
	IDEmisorFactura        string `xml:"sum1:IDEmisorFacturaAnulada"`
	NumSerieFactura        string `xml:"sum1:NumSerieFacturaAnulada"`
	FechaExpedicionFactura string `xml:"sum1:FechaExpedicionFacturaAnulada"`
}

// FacturaRectificada represents a rectified invoice
type FacturaRectificada struct {
	IDFactura IDFactura `xml:"sum1:IDFacturaRectificada"`
}

// FacturaSustituida represents a substituted invoice
type FacturaSustituida struct {
	IDFactura IDFactura `xml:"sum1:IDFacturaSustituida"`
}

// ImporteRectificacion contains rectification amounts
type ImporteRectificacion struct {
	BaseRectificada         num.Amount `xml:"sum1:BaseRectificada"`
	CuotaRectificada        num.Amount `xml:"sum1:CuotaRectificada"`
	CuotaRecargoRectificado num.Amount `xml:"sum1:CuotaRecargoRectificado,omitempty"`
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
	Impuesto                      string `xml:"sum1:Impuesto,omitempty"`
	ClaveRegimen                  string `xml:"sum1:ClaveRegimen,omitempty"`
	CalificacionOperacion         string `xml:"sum1:CalificacionOperacion,omitempty"`
	OperacionExenta               string `xml:"sum1:OperacionExenta,omitempty"`
	TipoImpositivo                string `xml:"sum1:TipoImpositivo,omitempty"`
	BaseImponibleOImporteNoSujeto string `xml:"sum1:BaseImponibleOimporteNoSujeto"`
	BaseImponibleACoste           string `xml:"sum1:BaseImponibleACoste,omitempty"`
	CuotaRepercutida              string `xml:"sum1:CuotaRepercutida,omitempty"`
	TipoRecargoEquivalencia       string `xml:"sum1:TipoRecargoEquivalencia,omitempty"`
	CuotaRecargoEquivalencia      string `xml:"sum1:CuotaRecargoEquivalencia,omitempty"`
}

// Encadenamiento contains chaining information between documents
type Encadenamiento struct {
	PrimerRegistro   string            `xml:"sum1:PrimerRegistro,omitempty"`
	RegistroAnterior *RegistroAnterior `xml:"sum1:RegistroAnterior,omitempty"`
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

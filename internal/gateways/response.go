package gateways

import (
	"github.com/nbio/xml"
)

// Envelope defines the SOAP envelope structure.
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	XMLNs   string   `xml:"xmlns:env,attr"`
	Body    Body     `xml:"Body"`
}

// Body represents the body of the SOAP envelope.
type Body struct {
	XMLName   xml.Name  `xml:"Body"`
	ID        string    `xml:"Id,attr,omitempty"`
	Fault     *Fault    `xml:"Fault,omitempty"`
	Respuesta *Response `xml:"RespuestaRegFactuSistemaFacturacion,omitempty"`
}

// Fault is issued by the SOAP server when something goes wrong.
type Fault struct {
	Code    string `xml:"faultcode"`
	Message string `xml:"faultstring"`
}

// Response defines the response fields from the VeriFactu gateway.
type Response struct {
	XMLName           xml.Name `xml:"RespuestaRegFactuSistemaFacturacion"`
	TikNamespace      string   `xml:"xmlns:tik,attr,omitempty"`
	TikRNamespace     string   `xml:"xmlns:tikR,attr,omitempty"`
	Cabecera          Cabecera `xml:"Cabecera"`
	TiempoEsperaEnvio int      `xml:"TiempoEsperaEnvio"`
	EstadoEnvio       string   `xml:"EstadoEnvio"`
	RespuestaLinea    []struct {
		IDFactura struct {
			IDEmisorFactura        string `xml:"IDEmisorFactura"`
			NumSerieFactura        string `xml:"NumSerieFactura"`
			FechaExpedicionFactura string `xml:"FechaExpedicionFactura"`
		} `xml:"IDFactura"`
		Operacion struct {
			TipoOperacion string `xml:"TipoOperacion"`
		} `xml:"Operacion"`
		EstadoRegistro           string `xml:"EstadoRegistro"`
		CodigoErrorRegistro      string `xml:"CodigoErrorRegistro,omitempty"`
		DescripcionErrorRegistro string `xml:"DescripcionErrorRegistro,omitempty"`
	} `xml:"RespuestaLinea"`
}

// Cabecera represents the header section of the response.
type Cabecera struct {
	ObligadoEmision struct {
		NombreRazon string `xml:"NombreRazon"`
		NIF         string `xml:"NIF"`
	} `xml:"ObligadoEmision"`
}

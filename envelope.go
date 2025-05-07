package verifactu

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"
)

// SUM is the namespace for the main VeriFactu schema
const (
	SUM          = "https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroLR.xsd"
	SUM1         = "https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroInformacion.xsd"
	EnvNamespace = "http://schemas.xmlsoap.org/soap/envelope/"
)

// for needed for timezones
var location *time.Location

// IssuerRole defines the role of the issuer in the invoice.
type IssuerRole string

// IssuerRole constants
const (
	IssuerRoleSupplier   IssuerRole = "E"
	IssuerRoleCustomer   IssuerRole = "D"
	IssuerRoleThirdParty IssuerRole = "T"
)

const (
	// CurrentVersion is the current version of the VeriFactu document
	CurrentVersion = "1.0"
)

func init() {
	var err error
	location, err = time.LoadLocation("Europe/Madrid")
	if err != nil {
		panic(err)
	}
}

// Envelope is the SOAP envelope wrapper used for sending messages to
// the remote service.
type Envelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	XMLNs   string   `xml:"xmlns:soapenv,attr"`
	SUM     string   `xml:"xmlns:sum,attr,omitempty"`
	SUM1    string   `xml:"xmlns:sum1,attr,omitempty"`
	Body    struct {
		ID             string          `xml:"soapenv:Id,attr,omitempty"`
		InvoiceRequest *InvoiceRequest `xml:"sum:RegFactuSistemaFacturacion,omitempty"`
	} `xml:"soapenv:Body"`
}

// EnvelopeResponse handles a SOAP response object that will correctly
// handle the namespaces.
type EnvelopeResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		ID              string           `xml:"Id,attr,omitempty"`
		Fault           *Fault           `xml:"Fault,omitempty"`
		InvoiceResponse *InvoiceResponse `xml:"RespuestaRegFactuSistemaFacturacion,omitempty"`
	} `xml:"Body"`
}

// Fault is issued by the SOAP server when something goes wrong.
type Fault struct {
	Code    string `xml:"faultcode"`
	Message string `xml:"faultstring"`
}

func newEnvelope() *Envelope {
	env := &Envelope{
		XMLNs: EnvNamespace,
		SUM:   SUM,
		SUM1:  SUM1,
	}
	return env
}

// Bytes returns the XML document bytes
func (d *Envelope) Bytes() ([]byte, error) {
	return toBytes(d)
}

// BytesIndent returns the indented XML document bytes
func (d *Envelope) BytesIndent() ([]byte, error) {
	return toBytesIndent(d)
}

func toBytes(doc any) ([]byte, error) {
	buf, err := buffer(doc, xml.Header, false)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func toBytesIndent(doc any) ([]byte, error) {
	buf, err := buffer(doc, xml.Header, true)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buffer(doc any, base string, indent bool) (*bytes.Buffer, error) {
	buf := bytes.NewBufferString(base)

	enc := xml.NewEncoder(buf)
	if indent {
		enc.Indent("", "  ")
	}

	if err := enc.Encode(doc); err != nil {
		return nil, fmt.Errorf("encoding document: %w", err)
	}

	return buf, nil
}

type timeLocationable interface {
	In(*time.Location) time.Time
}

func formatDateTimeZone(ts timeLocationable) string {
	return ts.In(location).Format("2006-01-02T15:04:05-07:00")
}

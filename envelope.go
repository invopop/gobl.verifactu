package verifactu

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// VeriFactu XML namespace constants
const (
	SUM          = "https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroLR.xsd"
	SUM1         = "https://www2.agenciatributaria.gob.es/static_files/common/internet/dep/aplicaciones/es/aeat/tike/cont/ws/SuministroInformacion.xsd"
	DS           = "http://www.w3.org/2000/09/xmldsig#"
	EnvNamespace = "http://schemas.xmlsoap.org/soap/envelope/"
)

// for needed for timezones
var location *time.Location

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
	DS      string   `xml:"xmlns:ds,attr,omitempty"`
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

// faultCodeRegexp is used to extract the VeriFactu error code embedded inside
// the fault string. Errors detected in the header of a request are always
// reported as SOAP faults, and the fault structure has no field available for
// the code, so the gateway includes it in the human readable message, e.g.
// `Codigo[4104].Error en la cabecera: ...`.
var faultCodeRegexp = regexp.MustCompile(`^Codigo\[(\d+)\]\.?\s*`)

// faultWhitespaceRegexp matches the runs of whitespace used by the gateway to
// indent the details appended to fault strings.
var faultWhitespaceRegexp = regexp.MustCompile(`\s+`)

// ErrorCode provides the VeriFactu error code of the fault, if one could be
// extracted from the fault string.
func (f *Fault) ErrorCode() string {
	m := faultCodeRegexp.FindStringSubmatch(f.faultString())
	if m == nil {
		return ""
	}
	return m[1]
}

// ErrorMessage provides the fault string without the error code prefix and
// with any indentation reduced to single spaces.
func (f *Fault) ErrorMessage() string {
	return strings.TrimSpace(faultCodeRegexp.ReplaceAllString(f.faultString(), ""))
}

// Err converts the fault into a structured error that the consumer can use to
// determine how the issue should be handled.
func (f *Fault) Err() *Error {
	e := ErrValidation
	if f.serverSide() {
		// Something went wrong inside the gateway, the request may be
		// worth sending again later.
		e = ErrServer
	}
	return e.WithCode(f.ErrorCode()).WithMessage(f.ErrorMessage())
}

// serverSide determines if the fault was caused by the remote service instead
// of the contents of our request. The prefix of the fault code is defined by
// the server, so only the local part can be compared.
func (f *Fault) serverSide() bool {
	code := f.Code
	if _, after, found := strings.Cut(code, ":"); found {
		code = after
	}
	return strings.EqualFold(code, "Server")
}

func (f *Fault) faultString() string {
	return strings.TrimSpace(faultWhitespaceRegexp.ReplaceAllString(f.Message, " "))
}

func newEnvelope() *Envelope {
	env := &Envelope{
		XMLNs: EnvNamespace,
		SUM:   SUM,
		SUM1:  SUM1,
		DS:    DS,
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

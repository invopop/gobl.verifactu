package verifactu

import "github.com/nbio/xml"

// InvoiceRequest represents the root element of a RegFactuSistemaFacturacion document
type InvoiceRequest struct {
	XMLName xml.Name              `xml:"sum:RegFactuSistemaFacturacion"`
	Header  *InvoiceRequestHeader `xml:"sum:Cabecera"`
	Lines   []*InvoiceRequestLine `xml:"sum:RegistroFactura,omitempty"`
}

// InvoiceRequestHeader contains the header information for a VeriFactu document
type InvoiceRequestHeader struct {
	Obligado           Issuer              `xml:"sum1:ObligadoEmision"`
	Representante      *Issuer             `xml:"sum1:Representante,omitempty"`
	RemisionVoluntaria *RemisionVoluntaria `xml:"sum1:RemisionVoluntaria,omitempty"`
	// RemisionRequerimiento *RemisionRequerimiento `xml:"sum1:RemisionRequerimiento,omitempty"` // not supported
}

// Issuer represents an obligated party in the document
type Issuer struct {
	NombreRazon string `xml:"sum1:NombreRazon"`
	NIF         string `xml:"sum1:NIF"`
}

// RemisionVoluntaria contains voluntary submission details
type RemisionVoluntaria struct {
	FechaFinVerifactu string `xml:"sum1:FechaFinVerifactu,omitempty"`
	Incidencia        string `xml:"sum1:Incidencia,omitempty"`
}

// InvoiceReuqestLine contains either an invoice registration or cancellation
type InvoiceRequestLine struct {
	Registration *InvoiceRegistration `xml:"sum1:RegistroAlta,omitempty"`
	Cancellation *InvoiceCancellation `xml:"sum1:RegistroAnulacion,omitempty"`
}

// AddRegistration adds the provided document to the list of registrations.
func (req *InvoiceRequest) AddRegistration(d *InvoiceRegistration) {
	d.NS = "" // Remove namespace
	req.addRow(&InvoiceRequestLine{
		Registration: d,
	})
}

// AddCancellation adds the requested document to the request body.
func (req *InvoiceRequest) AddCancellation(d *InvoiceCancellation) {
	req.addRow(&InvoiceRequestLine{
		Cancellation: d,
	})
}

func (req *InvoiceRequest) addRow(rf *InvoiceRequestLine) {
	if req.Lines == nil {
		req.Lines = make([]*InvoiceRequestLine, 0, 1)
	}
	req.Lines = append(req.Lines, rf)
}

// ChainData provides the chaining data for this line inside the
// invoice request.
func (line *InvoiceRequestLine) ChainData() *ChainData {
	if r := line.Registration; r != nil {
		return r.ChainData()
	}
	if r := line.Cancellation; r != nil {
		return r.ChainData()
	}
	return nil
}

// Envelop provides a SOAP Envelope around the InvoiceRequest, ready to
// send off via the API.
func (req *InvoiceRequest) Envelop() *Envelope {
	e := newEnvelope()
	e.Body.InvoiceRequest = req
	return e
}

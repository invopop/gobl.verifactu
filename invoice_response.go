package verifactu

import "github.com/nbio/xml"

// Default response status
const (
	StatusCorrect            string = "Correcto"
	StatusAcceptedWithErrors string = "AceptadoConErrores"
	StatusCancelled          string = "Anulada"
	StatusIncorrect          string = "Incorrecto"
)

// InvoiceResponse defines the response fields from the VeriFactu gateway.
type InvoiceResponse struct {
	XMLName xml.Name `xml:"RespuestaRegFactuSistemaFacturacion"`
	Header  struct {
		Issuer             InvoiceResponseIssuer               `xml:"ObligadoEmision"`
		Representative     *InvoiceResponseIssuer              `xml:"Representante,omitempty"`
		RemisionVoluntaria *InvoiceResponseVoluntarySubmission `xml:"sum1:RemisionVoluntaria,omitempty"`
	} `xml:"Cabecera"`
	Wait   int                    `xml:"TiempoEsperaEnvio"`
	Status string                 `xml:"EstadoEnvio"`
	Lines  []*InvoiceResponseLine `xml:"RespuestaLinea"`
}

// InvoiceResponseIssuer maps the response from the invoice request
// for the issuer.
type InvoiceResponseIssuer struct {
	Name string `xml:"NombreRazon"`
	NIF  string `xml:"NIF"`
}

// InvoiceResponseVoluntarySubmission for additional header data.
type InvoiceResponseVoluntarySubmission struct {
	Date     string `xml:"FechaFinVerifactu,omitempty"`
	Incident string `xml:"Incidencia,omitempty"`
}

// InvoiceResponseLine defines the contents of a single line of the request.
type InvoiceResponseLine struct {
	XMLName xml.Name `xml:"RespuestaLinea"`
	ID      struct {
		Issuer string `xml:"IDEmisorFactura"`
		Code   string `xml:"NumSerieFactura"`
		Date   string `xml:"FechaExpedicionFactura"`
	} `xml:"IDFactura"`
	Operation struct {
		Type             string `xml:"TipoOperacion"`
		Correction       string `xml:"Subsanacion,omitempty"`
		RejectedPrevious string `xml:"RechazoPrevio,omitempty"`
		NoPrevious       string `xml:"SinRegistroPrevio,omitempty"`
	} `xml:"Operacion"`
	Ref         string `xml:"RefExterna,omitempty"`
	Status      string `xml:"EstadoRegistro"`
	Code        string `xml:"CodigoErrorRegistro,omitempty"`
	Description string `xml:"DescripcionErrorRegistro,omitempty"`

	// Duplicated contains information about possible duplicated requests
	// that may need to be handled differently.
	Duplicated *InvoiceResponseLineDuplicated `xml:"RegistroDuplicado,omitempty"`
}

// InvoiceResponseLineDuplicated describes details about the detected duplicate.
type InvoiceResponseLineDuplicated struct {
	ID          string `xml:"IdPeticionRegistroDuplicado"`
	Status      string `xml:"EstadoRegistroDuplicado"`
	Code        string `xml:"CodigoErrorRegistro,omitempty"`
	Description string `xml:"DescripcionErrorRegistro,omitempty"`
}

// Message provides a message body, if any.
func (r *InvoiceResponseLine) Message() string {
	txt := r.Description
	if txt != "" {
		return r.Status + ": " + txt
	}
	return r.Status
}

// Accepted will return true if the line was accepted.
func (r *InvoiceResponseLine) Accepted() bool {
	switch r.Status {
	case StatusCorrect, StatusCancelled, StatusAcceptedWithErrors:
		return true
	}
	return false
}

// Warning will return a warning error if the line was accepted, but
// contains additional details that should be investigated.
func (r *InvoiceResponseLine) Warning() error {
	if r.Status == StatusAcceptedWithErrors {
		return ErrWarning.WithCode(r.Code).WithMessage(r.Message())
	}
	return nil
}

// Bytes prepares an indendented XML document suitable for persistence.
func (ir *InvoiceResponse) Bytes() ([]byte, error) {
	return toBytesIndent(ir)
}

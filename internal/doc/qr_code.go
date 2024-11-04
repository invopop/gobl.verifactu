package doc

import (
	"fmt"
	"net/url"
	// "github.com/sigurn/crc8"
)

// Codes contain info about the codes that should be generated and shown on a
// Ticketbai invoice. One is an alphanumeric code that identifies the invoice
// and the other one is a URL (which can be used by a customer to validate that
// the invoice has been sent to the tax agency) that should be encoded as a
// QR code in the printed invoice / ticket.
type Codes struct {
	URLCode string
	QRCode  string
}

const (
	BaseURL = "https://prewww2.aeat.es/wlpl/TIKE-CONT/ValidarQR?"
)

// var crcTable = crc8.MakeTable(crc8.CRC8)

// generateCodes will generate the QR and URL codes for the invoice
func (doc *VeriFactu) generateCodes(inv *RegistroAlta) *Codes {
	urlCode := doc.generateURLCode(inv)
	// qrCode := doc.generateQRCode(urlCode)

	return &Codes{
		URLCode: urlCode,
		// QRCode:  qrCode,
	}
}

// generateURLCode generates the encoded URL code with parameters.
func (doc *VeriFactu) generateURLCode(inv *RegistroAlta) string {
	// URL encode each parameter
	nif := url.QueryEscape(doc.Cabecera.Obligado.NIF)
	numSerie := url.QueryEscape(inv.IDFactura.NumSerieFactura)
	fecha := url.QueryEscape(inv.IDFactura.FechaExpedicionFactura)
	importe := url.QueryEscape(fmt.Sprintf("%.2f", inv.ImporteTotal))

	// Build the URL
	urlCode := fmt.Sprintf("%snif=%s&numserie=%s&fecha=%s&importe=%s",
		BaseURL, nif, numSerie, fecha, importe)

	return urlCode
}

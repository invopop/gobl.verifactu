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
	TestURL = "https://prewww2.aeat.es/wlpl/TIKE-CONT/ValidarQR?"
	ProdURL = "https://www2.agenciatributaria.gob.es/wlpl/TIKE-CONT/ValidarQR?nif=89890001K&numserie=12345678-G33&fecha=01-09-2024&importe=241.4"
)

// var crcTable = crc8.MakeTable(crc8.CRC8)

// generateCodes will generate the QR and URL codes for the invoice
func (doc *VeriFactu) generateCodes() *Codes {
	urlCode := doc.generateURLCodeAlta()
	// qrCode := doc.generateQRCode(urlCode)

	return &Codes{
		URLCode: urlCode,
		// QRCode:  qrCode,
	}
}

// generateURLCode generates the encoded URL code with parameters.
func (doc *VeriFactu) generateURLCodeAlta() string {

	// URL encode each parameter
	nif := url.QueryEscape(doc.RegistroFactura.RegistroAlta.IDFactura.IDEmisorFactura)
	numSerie := url.QueryEscape(doc.RegistroFactura.RegistroAlta.IDFactura.NumSerieFactura)
	fecha := url.QueryEscape(doc.RegistroFactura.RegistroAlta.IDFactura.FechaExpedicionFactura)
	importe := url.QueryEscape(doc.RegistroFactura.RegistroAlta.ImporteTotal)

	// Build the URL
	urlCode := fmt.Sprintf("%snif=%s&numserie=%s&fecha=%s&importe=%s",
		TestURL, nif, numSerie, fecha, importe)

	return urlCode
}

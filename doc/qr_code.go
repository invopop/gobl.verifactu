package doc

import (
	"fmt"
	"net/url"
)

const (
	testURL = "https://prewww2.aeat.es/wlpl/TIKE-CONT/ValidarQR?"
	prodURL = "https://www2.agenciatributaria.gob.es/wlpl/TIKE-CONT/ValidarQR?"
)

// generateURL generates the encoded URL code with parameters.
func (doc *Envelope) generateURL(production bool) string {
	nif := url.QueryEscape(doc.Body.VeriFactu.RegistroFactura.RegistroAlta.IDFactura.IDEmisorFactura)
	numSerie := url.QueryEscape(doc.Body.VeriFactu.RegistroFactura.RegistroAlta.IDFactura.NumSerieFactura)
	fecha := url.QueryEscape(doc.Body.VeriFactu.RegistroFactura.RegistroAlta.IDFactura.FechaExpedicionFactura)
	importe := url.QueryEscape(doc.Body.VeriFactu.RegistroFactura.RegistroAlta.ImporteTotal)

	if production {
		return fmt.Sprintf("%s&nif=%s&numserie=%s&fecha=%s&importe=%s", prodURL, nif, numSerie, fecha, importe)
	}
	return fmt.Sprintf("%snif=%s&numserie=%s&fecha=%s&importe=%s", testURL, nif, numSerie, fecha, importe)
}

package verifactu

import (
	"fmt"
	"net/url"
)

const (
	testURL = "https://prewww2.aeat.es/wlpl/TIKE-CONT/ValidarQR?"
	prodURL = "https://www2.agenciatributaria.gob.es/wlpl/TIKE-CONT/ValidarQR?"
)

// generateURL generates the encoded URL code with parameters.
func (r *InvoiceRegistration) generateURL(production bool) string {
	nif := url.QueryEscape(r.IDFactura.IDEmisorFactura)
	numSerie := url.QueryEscape(r.IDFactura.NumSerieFactura)
	fecha := url.QueryEscape(r.IDFactura.FechaExpedicionFactura)
	importe := url.QueryEscape(r.ImporteTotal)

	if production {
		return fmt.Sprintf("%s&nif=%s&numserie=%s&fecha=%s&importe=%s", prodURL, nif, numSerie, fecha, importe)
	}
	return fmt.Sprintf("%snif=%s&numserie=%s&fecha=%s&importe=%s", testURL, nif, numSerie, fecha, importe)
}
